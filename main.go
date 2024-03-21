package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/domicon-labs/kzg-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
)

const dChunkSize = 30
const dSrsSize = 1 << 16

type stDA struct {
	index       uint64
	legth       uint64
	broadcaster common.Address
	user        common.Address
	commitment  hexutil.Bytes
	sign        hexutil.Bytes
	data        hexutil.Bytes
}

func generateDA() *stDA {
	sdk, err := kzg_sdk.InitDomiconSdk(dSrsSize, "./srs")
	if err != nil {
		log.Println("InitDomiconSdk failed")
		return nil
	}

	fmt.Println("2. prepare test data ")
	data := make([]byte, dChunkSize*17)
	for i := range data {
		data[i] = 1
	}

	log.Println("3. generate data commit")
	length := len(data)
	log.Println("len", length)
	digest, err := sdk.GenerateDataCommit(data)
	if err != nil {
		log.Println("GenerateDataCommit failed")
		return nil
	}
	log.Println("commit data is:", digest.Bytes())
	CM := digest.Bytes()
	log.Println("cm:", hex.EncodeToString(CM[:]))
	singer := kzg_sdk.NewEIP155FdSigner(big.NewInt(31337))

	// key, _ := crypto.GenerateKey()
	pkBytes, err := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Println("DecodeString failed")
		return nil
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
	log.Println("key:", privateKeyString)

	log.Println("senAddr----", senAddr.Hex())

	index := 1
	_, sigData, err := kzg_sdk.SignFd(senAddr, senAddr, 0, uint64(index), uint64(length), CM[:], singer, privateKey)
	if err != nil {
		log.Println("msg", "err", err)
		return nil
	}
	log.Println("sigdata:", hex.EncodeToString(sigData))
	log.Println("len:", len(sigData))
	return &stDA{
		index:       uint64(index),
		legth:       uint64(length),
		broadcaster: senAddr,
		user:        senAddr,
		commitment:  hexutil.Bytes(CM[:]),
		sign:        hexutil.Bytes(sigData),
		data:        hexutil.Bytes(data),
	}
}

func main() {
	url := "http://127.0.0.1:8200"
	rpcCli, err := rpc.DialOptions(context.Background(), url)
	if err != nil {
		log.Println("msg", "err", err)
		return
	} else {
		log.Println("success")
	}
	rawDA := generateDA()
	var result *[]byte
	err = rpcCli.CallContext(context.Background(), result, "sendDA", rawDA.index, rawDA.legth, rawDA.broadcaster,
		rawDA.user, rawDA.commitment, rawDA.sign, rawDA.data)
	if err != nil {
		log.Println("msg", "err", err)
	} else {
		log.Println("msg", "respose", result)
	}
}
