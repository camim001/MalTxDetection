const wallets= [
                      '0x214e4dAba60AEDA376B77b4779D5c645358b6848', //3
                      '0xFD84B7a430873dF2761C5557eB5650FF6Ad2Abae', //3
                      '0xF4afa0480220054433decE855d1d2aCFaada6e97', //3
                      '0x67D97506a68ee42a42A39D957E2fC39C288c0033', //2
                      '0x155B687CFdFCe8617029445f5B4C71901F58f0C4', //2
                      '0x485B5CcA6529C573a5C8080A5bE46b994970250D', //2
                      '0xc061a63733CfC820A6758c3CFB866C2a630a09D1', //1
                      '0xA5f2210104d4839Ea9e5fAC14DC73Ef5eB02e9aa', //1
                      '0xCAee0d3e83b4Ce5da15De7Ba6E9D40184F8dD61E'  //1
              ]

const delay = require('delay');
const Web3 = require('web3');
const infuraUrl01 = 'http://localhost:8101';
const infuraUrl02 = 'http://localhost:8102';
const infuraUrl03 = 'http://localhost:8103';
const web3_01 = new Web3(infuraUrl01);
const web3_02 = new Web3(infuraUrl02);
const web3_03 = new Web3(infuraUrl03);

myFunc()

async function myFunc(){

  for(i=0;i<29;i++){
  
  		lpSTRT = Date.now()
  		console.log("Start of Loop ",i,"fromAddr",",","toAddr",",",lpSTRT,",","transactionHash",",","loop")
	
		await sendTx(wallets[0],wallets[3],i,web3_03);
		await sendTx(wallets[0],wallets[5],i,web3_03);
		await sendTx(wallets[0],wallets[6],i,web3_03);
		await sendTx(wallets[0],wallets[8],i,web3_03);

		await sendTx(wallets[1],wallets[3],i,web3_03);
		await sendTx(wallets[1],wallets[5],i,web3_03);
		await sendTx(wallets[1],wallets[6],i,web3_03);
		await sendTx(wallets[1],wallets[8],i,web3_03);

		await sendTx(wallets[2],wallets[3],i,web3_03);
		await sendTx(wallets[2],wallets[5],i,web3_03);
		await sendTx(wallets[2],wallets[6],i,web3_03);
		await sendTx(wallets[2],wallets[8],i,web3_03);

		await sendTx(wallets[3],wallets[0],i,web3_02);
		await sendTx(wallets[3],wallets[2],i,web3_02);
		await sendTx(wallets[3],wallets[6],i,web3_02);
		await sendTx(wallets[3],wallets[8],i,web3_02);

		await sendTx(wallets[4],wallets[0],i,web3_02);
		await sendTx(wallets[4],wallets[2],i,web3_02);
		await sendTx(wallets[4],wallets[6],i,web3_02);
		await sendTx(wallets[4],wallets[8],i,web3_02);

		await sendTx(wallets[5],wallets[0],i,web3_02);
		await sendTx(wallets[5],wallets[2],i,web3_02);
		await sendTx(wallets[5],wallets[6],i,web3_02);
		await sendTx(wallets[5],wallets[8],i,web3_02);

		await sendTx(wallets[6],wallets[0],i,web3_01);
		await sendTx(wallets[6],wallets[2],i,web3_01);
		await sendTx(wallets[6],wallets[3],i,web3_01);
		await sendTx(wallets[6],wallets[5],i,web3_01);

		await sendTx(wallets[7],wallets[0],i,web3_01);
		await sendTx(wallets[7],wallets[2],i,web3_01);
		await sendTx(wallets[7],wallets[3],i,web3_01);
		await sendTx(wallets[7],wallets[5],i,web3_01);

		await sendTx(wallets[8],wallets[0],i,web3_01);
		await sendTx(wallets[8],wallets[2],i,web3_01);
		await sendTx(wallets[8],wallets[3],i,web3_01);
		await sendTx(wallets[8],wallets[5],i,web3_01);
		console.log("End of Loop ",i,"fromAddr",",","toAddr",",",(Date.now() - lpSTRT),",","transactionHash",",","loop")
	}
}

async function sendTx(fromAddr,toAddr,loop,web3){

	try{
	
		stRT = Date.now()
		await web3.eth.personal.unlockAccount(fromAddr, "12345")
		aa = await web3.eth.sendTransaction({
	      		from:fromAddr,
      			to: toAddr,
      			value: '1'
    		});
    		console.log("Success",",",fromAddr,",",toAddr,",",(Date.now() - stRT),",",aa.transactionHash,",",loop)

  	}catch(e){
		
		console.log("Fail",",",fromAddr,",",toAddr,",",(Date.now() - stRT),",","Failed Transaction",",",loop)
  	}
	
}
