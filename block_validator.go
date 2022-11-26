// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/common"

	// these are requied for the transaction checking
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	// these are required for transaction checking


)

// BlockValidator is responsible for validating block headers, uncles and
// processed state.
//
// BlockValidator implements Validator.
type BlockValidator struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}


//******************************************************************
// Define variables that will be used by the checking function
// which uses a web service to check
//******************************************************************
var vlddbUrl string //define the default IP where the server is
var vldfilePath string //the file path

//Define the record type that will be received from teh server
type vldaddrRecord struct{
	Status string
	Time string
	Reason string
}

//******************************************************************
// Define variables that will be used by the checking function
// which uses a local data base
//******************************************************************
type vldmalAdr struct{
  Value string
  Status string
  Time string
  Reason string
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody validates the given block's uncles and verifies the block
// header's transaction and uncle roots. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(block *types.Block) error {

//*******************************************************************
// We can validate for malicious transactions here
// (*(*block.transactions[0]).from.v.(data).from   -> from address
// (*(*(*block).transactions[0]).inner.(data)).To   -> To address



	var txArray = block.Transactions()
	fmt.Printf("Transaction Array is: %T\n", txArray)
	//fmt.Println(gorg[0].Hash())

	if txArray.Len() != 0 {

		fmt.Println("*********************************************************************")
		fmt.Println("Check the transactions in mined blocks")
		fmt.Println("*********************************************************************")

		for token, worth  := range txArray {
			fmt.Println(token, worth)
			//fmt.Println(worth.ChainId())

			// extract the sender of teh transaction
			signer := types.NewEIP155Signer(worth.ChainId())
			sender, _ := signer.Sender(worth)

			//if vldcheckLocal(sender){
				//fmt.Println("Malicious Sender")
				//sendWarning(sender)
			//}

			if vldbeforeCheckAddress(sender){
				fmt.Println("Sender is Malicious")
				//here we need to do a routine to inform another web service
				//of the malicious address
				sendWarning(sender)
			}

//			if vldcheckLocal(*worth.To()){
//				fmt.Println("Malicious Receiver")
				//send malicious receiver
//				sendWarning(*worth.To())
//			}

			if vldbeforeCheckAddress(*worth.To()){
				fmt.Println("Receiver is Malicious")
				//here we need to do a routine to inform another web service
				//of the malicious address
				sendWarning(*worth.To())
			}
			//fmt.Printf("type of the sender is: %T\n", sender.Hex())
			//fmt.Printf("sender: %v\n", sender.Hex())
		  //var fromAddress = common.HexToAddress(sender.Hex())
			//fmt.Printf("type of the sender is: %T\n", fromAddress)
			//extract the sender of the transaction

			//fmt.Printf("The value is of type : %T\n", worth)
			//fmt.Printf("Source Tx is of type: %T\n", worth.To())
			//fmt.Println(worth.To())

		}

	} else {
		fmt.Println("Transaction Array is empty")
	}
//*******************************************************************


	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		return ErrKnownBlock
	}
	// Header validity is known at this point, check the uncles and transactions
	header := block.Header()
	if err := v.engine.VerifyUncles(v.bc, block); err != nil {
		return err
	}
	if hash := types.CalcUncleHash(block.Uncles()); hash != header.UncleHash {
		return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	}
	if hash := types.DeriveSha(block.Transactions(), trie.NewStackTrie(nil)); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}
	if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			return consensus.ErrUnknownAncestor
		}
		return consensus.ErrPrunedAncestor
	}
	return nil

}

// ValidateState validates the various changes that happen after a state
// transition, such as amount of used gas, the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(block *types.Block, statedb *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	header := block.Header()
	if block.GasUsed() != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", block.GasUsed(), usedGas)
	}
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, Rn]]))
	receiptSha := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}

	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(v.config.IsEIP158(header.Number)); header.Root != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	return nil
}

// CalcGasLimit computes the gas limit of the next block after parent. It aims
// to keep the baseline gas close to the provided target, and increase it towards
// the target if the baseline gas is lower.
func CalcGasLimit(parentGasLimit, desiredLimit uint64) uint64 {
	delta := parentGasLimit/params.GasLimitBoundDivisor - 1
	limit := parentGasLimit
	if desiredLimit < params.MinGasLimit {
		desiredLimit = params.MinGasLimit
	}
	// If we're outside our allowed gas range, we try to hone towards them
	if limit < desiredLimit {
		limit = parentGasLimit + delta
		if limit > desiredLimit {
			limit = desiredLimit
		}
		return limit
	}
	if limit > desiredLimit {
		limit = parentGasLimit - delta
		if limit < desiredLimit {
			limit = desiredLimit
		}
	}
	return limit
}


//**************************************************************************************************
//     CHECK THE Addresses
//**************************************************************************************************
func vldbeforeCheckAddress( malAddress common.Address) bool{

	// declare the variables
	var strMalAddr string
	var connectUrl string
	var checkRecord addrRecord				// the record type

	//vlddbUrl = "http://192.168.1.201:8081/"
	vlddbUrl = "http://134.122.102.42:8080/"
	// convert the address to string
	strMalAddr = malAddress.Hex()
	// construct the URL with address
	connectUrl = vlddbUrl+strMalAddr

	resp, err := http.Get(connectUrl)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		 fmt.Println(err)
	}
	//Convert the body to type string
	sb := string(body)
	fmt.Println(sb)

	//convert the string which is json to a structure
	json.Unmarshal([]byte(sb), &checkRecord)

	if(checkRecord.Status == "malicious"){
		resp.Body.Close()
		return true
	}

	resp.Body.Close()
	return false

}

//*******************************************************************************************
// Check the addresses from a Local file
//*******************************************************************************************
func vldcheckLocal(malAddress common.Address) bool {

	var strMalAddr string
	var malAdrArray []vldmalAdr  //create an array for the malicious addresses
	var hazin bool

	vldfilePath = "/home/martin/go/myPrivate/malAddr.json"

	// convert the address to string
	strMalAddr = malAddress.Hex()

	jsonFile, err := os.Open(vldfilePath)

	// if we os.Open returns an error then handle it
	if err != nil {
	    fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


 	byteValue, _ := ioutil.ReadAll(jsonFile)

	//convert the string which is json to a structure
	json.Unmarshal([]byte(byteValue), &malAdrArray)

	for _,s := range malAdrArray{
		if (s.Value == strMalAddr){
			hazin = true
		}
	}
	jsonFile.Close()

	if(hazin){
		return true
	}else{
		return false
	}

}

//*************************************************************************
// Send the warning to a remote server
//*************************************************************************
func sendWarning(malAddress common.Address){

	var strMalAddr string
	var connectUrl string
	var vldWarnUrl string

	//vldWarnUrl = "http://192.168.1.201:8081/martin/"
	vldWarnUrl = "http://134.122.102.42:8080/martin/"

	strMalAddr = malAddress.Hex()
	connectUrl = vldWarnUrl+strMalAddr+"%20is%20malicious%0D%0A"

	resp, err := http.Get(connectUrl)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(resp)
	}

}
