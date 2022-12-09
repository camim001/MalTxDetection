# MalTxDetection
Ethereum Miner Code modified to include malicious transaction detetcion

This reposirtory contains the following:

Markup: * Bullet list
  * The modified ETH source code to enable malicious transaction detection and prevention during transaction insterion into the transaction pool
    *  tx_poolSanitized.go
  * The modified GETH soruce code to enable malicious transaction detection during block verification
    *   block_validatorSanitized.go
  * Test scripts used for automated transaction processing in JS
    * SendTransactionsTest.js
    * SendTransactionsNoAsynchV4.js
  * Samples of output log files from automated test scripts 
    * AllForCENT.txt
    * MixNodeCENT.txt
  * Log files imported into Excel
    * AllForCENT.xlsx
    * MixNodeCENT.xlsx
