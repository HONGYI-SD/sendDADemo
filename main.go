package main

import (
	"fmt"
	"math/big"

	kzg_sdk "github.com/domicon-labs/kzg-sdk"
	"github.com/ethereum/go-ethereum/crypto"
)

const dChunkSize = 30
const dSrsSize = 1 << 16

func main() {
	sdk, err := kzg_sdk.InitDomiconSdk(dSrsSize, "./srs")
	if err != nil {
		fmt.Println("InitDomiconSdk failed")
		return
	}

	fmt.Println("2. prepare test data ")
	data := make([]byte, dChunkSize*17)
	for i := range data {
		data[i] = 1
	}

	fmt.Print("3. generate data commit")
	length := len(data)
	digest, err := sdk.GenerateDataCommit(data)
	if err != nil {
		fmt.Println("GenerateDataCommit failed")
		return
	}
	fmt.Println("commit data is:", digest.Bytes())
	CM := digest.Bytes()
	singer := kzg_sdk.NewEIP155FdSigner(big.NewInt(31337))

	key, _ := crypto.GenerateKey()
	senAddr := crypto.PubkeyToAddress(key.PublicKey)
	privateKeyString := fmt.Sprintf("Private Key: %x", key.D)
	fmt.Println("key:", privateKeyString)

	println("senAddr----", senAddr.Hex())

	index := 1
	sigHash, sigData, err := kzg_sdk.SignFd(senAddr, senAddr, 0, uint64(index), uint64(length), CM[:], singer, key)
	if err != nil {
		return
	}
	fmt.Println("sigHash:", sigHash.Hex())
	fmt.Println("len:", len(sigData))
}
