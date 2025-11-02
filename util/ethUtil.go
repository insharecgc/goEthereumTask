package util

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SendNewTransaction(
	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	to common.Address,
	value *big.Int,
	gasLimit uint64,
	data []byte,
) (*types.Transaction, error) {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("fromAddress:", fromAddress.Hex())
	// 获取当前账户的nonce值
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	// 获取建议gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	// 获取链ID
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("chainId:", chainId.String(), "nonce:", nonce, "gasPrice:", gasPrice.String())

	// 1.创建一个新的交易对象
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	// 2.使用私钥对交易进行签名
	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	// 3.验证签名是否正确
	recoveredAddr, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}
	if recoveredAddr != fromAddress {
		return nil, fmt.Errorf("签名验证失败，签名地址与发送地址不匹配")
	}

	// 4.发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}
