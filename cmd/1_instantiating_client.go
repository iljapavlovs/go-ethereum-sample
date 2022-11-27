package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func main() {

	/** 0. ABOUT TRANSPORT PROTOCOL - TO WHAT TO CONNECT TO
	There are three transport protocols that can be used to connect the Javascript environment to Geth:

	1. IPC (Inter-Process Communication): Provides unrestricted access to all APIs, but only works when the
	console is run on the same host as the Geth node.

	2. HTTP: By default provides access to the eth, web3 and net method namespaces.

	3. Websocket: By default provides access to the eth, web3 and net method namespaces.
	*/

	/* 1. CREATE CLIENT TO THE RPC JSON NODE
	The client is an instance of the Client struct which has associated functions that wrap requests to the Ethereum or Geth RPC API endpoints.

	A client is instantiated by passing a raw url or path to an ipc file to the client’s Dial function. In the following code snippet the path to the ipc file for a local Geth node is provided to ethclient.Dial()
	*/

	// create instance of ethclient and assign to client
	client, err := ethclient.Dial("https://eth-goerli.g.alchemy.com/v2/Og_z2jV9M75kZSYYEolGKrm2LmFcGAR3")
	if err != nil {
		panic(err)
	}

	/*
		2. USE CLIENT TO INTERACT WITH JSON RPC API
			The client can now be used to handle requests to the Geth node using the full JSON-RPC API. For example,
			* the function BlockNumer() wraps a call to the eth_blockNumber endpoint.
			* The function SendTransaction wraps a call to eth_sendTransaction
	*/

	/*
		3. CONTEXT DEFINES CONTEXT ABOUT REQUESTS SENT FROM THE APP (DEADLINES, CANCELLATION SIGNALS) (GOLANG FEATURE)
			Context.Background() - CREATE EMPTY CONTEXT
	*/

	//GET CHAIN ID - e.g. is needed when signing a transaction
	chainid, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ChainID: %d\n", chainid)

	//GET NONCE AT SPECIFIC BLOCK NUMBER
	addr := common.HexToAddress("0x23b5613fc04949F4A53d1cc8d6BCCD21ffc38C11")
	nonce, err := client.NonceAt(context.Background(), addr, big.NewInt(7916447))
	// check https://goerli.etherscan.io/tx/0x3da8d3d995298dfd2fe288ac82e14248e6714f20fd3684873e954b5ebdadb530
	fmt.Printf("NonceAt: %d Nonce for 0x23b5613fc04949F4A53d1cc8d6BCCD21ffc38C11 Address for 7965747 blocknumber\n", nonce)

	pendingNonce, err := client.PendingNonceAt(context.Background(), addr)
	fmt.Println("PendingNonce: ", pendingNonce)

	/*
		4. Querying past events - GET LOGS FROM THE PAST AND APPLY A FILTER
	*/

	blockNum, err := client.BlockNumber(context.Background()) // returns the most recent block number
	if err != nil {
		panic(err)
	}

	fmt.Println("BlockNumber: ", blockNum)

	//blockNumConverted := new(big.Int).SetUint64(blockNum)
	//q := ethereum.FilterQuery{
	//	FromBlock: new(big.Int).Sub(blockNumConverted, big.NewInt(10)),
	//	ToBlock:   blockNumConverted,
	//	Topics:    [][]common.Hash{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
	//}
	//logs, err := client.FilterLogs(context.Background(), q)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("FilterLogs: ", logs)

	err = sendTransaction(client)
	if err != nil {
		panic(err)
	}

}

// sendTransaction sends a transaction with 1 ETH to a specified address.
func sendTransaction(cl *ethclient.Client) error {
	var (
		SK       = "10c2b3b0f10b4ef0086854a66c78dab4ea31aa5624af83b7ccd99012790aac6a"
		ADDR     = "0x23b5613fc04949F4A53d1cc8d6BCCD21ffc38C11"
		sk       = crypto.ToECDSAUnsafe(common.FromHex(SK))
		to       = common.HexToAddress("0x975Cbd9C3c5863a82484cF0278E5C21f2aE9761b")
		value    = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
		sender   = common.HexToAddress(ADDR)
		gasLimit = uint64(21000)
	)
	// Retrieve the chainid (needed for signer)
	chainid, err := cl.ChainID(context.Background())
	if err != nil {
		return err
	}
	// Retrieve the pending nonce
	nonce, err := cl.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return err
	}
	// Get suggested gas price
	tipCap, _ := cl.SuggestGasTipCap(context.Background())
	feeCap, _ := cl.SuggestGasPrice(context.Background())
	// Create a new transaction
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID: chainid,
			Nonce:   nonce,
			//TODO - not necessary to multiply by 2
			GasTipCap: tipCap.Mul(tipCap, big.NewInt(2)),
			GasFeeCap: feeCap.Mul(tipCap, big.NewInt(2)),
			Gas:       gasLimit,
			To:        &to,
			Value:     value,
			Data:      nil,
		})
	// Sign the transaction using our keys
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainid), sk)
	if err != nil {
		panic(err)
	}

	// Send the transaction to our node
	err = cl.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	receipt, err := bind.WaitMined(context.Background(), cl, signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Println("TX Receipt: ", *receipt)
	return err
}
