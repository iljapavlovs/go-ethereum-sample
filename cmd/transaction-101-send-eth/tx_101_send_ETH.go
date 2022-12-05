package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"log"
	"math/big"
)

func main() {
	// 1. CONNECT TO THE JSON RPC NODE CLIENT
	client, err := ethclient.Dial("https://eth-goerli.g.alchemy.com/v2/Og_z2jV9M75kZSYYEolGKrm2LmFcGAR3")
	if err != nil {
		log.Fatal(err)
	}

	// 2. LOAD PRIVATE KEY
	privateKey, err := crypto.HexToECDSA("10c2b3b0f10b4ef0086854a66c78dab4ea31aa5624af83b7ccd99012790aac6a")
	if err != nil {
		log.Fatal(err)
	}

	/*
		Every transaction requires a nonce. A nonce by definition is a number that is only used once. If it's a new account sending out a transaction then the nonce will be 0. Every new transaction from an account must have a nonce that the previous nonce incremented by 1. It's hard to keep manual track of all the nonces so the ethereum client provides a helper method PendingNonceAt that will return the next nonce you should use.

		The function requires the public address of the account we're sending from -- which we can derive from the private key
	*/

	// 3. GET PUBLIC KEY FROM PRIVATE KEY
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 4. GET NONCE FOR THE TX
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 5. SET AMOUNT TO SEND
	value := new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)) // is the same as big.NewInt(1000000000000000000) // in wei (1 eth)

	// 6. SET GAS LIMIT (standard ETH transfer is 21000 units.)
	gasLimit := uint64(21000) // in units

	// 7. SET GAS PRICE you are willing to pay in order to get your tx asap to blockchain
	//gasPrice := big.NewInt(30000000000) // in wei (30 gwei) - HARDCODING GAS PRICE - NOT RECOMMENDED

	//go-ethereum provides "SuggestGasPrice" function for getting the average gas price based on x number of previous blocks
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 7.1 WHEN USING EIP1559 TXS - need to set MAX PRIORITY FEE
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 8. GET RECEPIANT ADDRESS
	toAddress := common.HexToAddress("0x4ed03F492CeD0487eEA9fd93eDf04A6896A83CC8")

	// 9. GET CHAIN ID
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 10. CREATE UNSIGNED TX
	//tx := types.NewTx(nonce, toAddress, value, gasLimit, gasPrice, nil)

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     nonce,
			GasTipCap: gasTipCap,
			GasFeeCap: gasPrice,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     value,
			Data:      nil,
		})

	// 11. SIGN TX WITH PRIVATE KEY OF THE SENDER
	//  using NewLondonSigner signer since it returns all available signers, also could use NewEIP155Signer
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 12. BROADCAST TX TO ENTIRE BLOCKCHAIN
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 13. WAIT UNTIL TX IS MINED ON THE BLOCKCHAIN
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TX Receipt: ", *receipt)
	fmt.Printf("TX https://goerli.etherscan.io/tx/%s", receipt.TxHash)
}
