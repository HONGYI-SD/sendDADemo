package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1"
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

	fmt.Println("3. generate data commit")
	length := len(data)
	fmt.Println("len", length)
	digest, err := sdk.GenerateDataCommit(data)
	if err != nil {
		fmt.Println("GenerateDataCommit failed")
		return
	}
	fmt.Println("commit data is:", digest.Bytes())
	CM := digest.Bytes()
	fmt.Println("cm:", hex.EncodeToString(CM[:]))
	singer := kzg_sdk.NewEIP155FdSigner(big.NewInt(31337))

	// key, _ := crypto.GenerateKey()
	pkBytes, err := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		fmt.Println("DecodeString failed")
		return
	}
	privateKeyInt := new(big.Int).SetBytes(pkBytes)
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(), // 曲线类型，此处使用P-256曲线
			X:     nil,              // X和Y值需要通过私钥计算，所以设置为nil
			Y:     nil,
		},
		D: privateKeyInt, // 私钥值
	}
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	senAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyString := fmt.Sprintf("Private Key: %x", privateKey.D)
	fmt.Println("key:", privateKeyString)

	println("senAddr----", senAddr.Hex())

	index := 1
	_, sigData, err := kzg_sdk.SignFd(senAddr, senAddr, 0, uint64(index), uint64(length), CM[:], singer, privateKey)
	if err != nil {
		return
	}
	fmt.Println("sigdata:", hex.EncodeToString(sigData))
	fmt.Println("len:", len(sigData))
}
