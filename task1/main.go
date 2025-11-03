package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/insharecgc/goEthereumTask/internal/contractbindings"
	"github.com/insharecgc/goEthereumTask/internal/util"
)

const infuraKey = "d56339e6f6a0412ea3ae3710fe72f198"

var testPrivateKey, _ = crypto.HexToECDSA("your secret key")

// 指定区块号查询区块信息
func queryBlockInfo(blockNumber *big.Int) {
	client, _ := util.GetHttpClient(infuraKey)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blcok hash:", block.Hash().Hex())
	fmt.Println("blcok time:", block.Time())
	fmt.Println("blcok tx len:", block.Transactions().Len())
	fmt.Println("blcok size:", block.Size())
}

// 发送交易
func transferETH(privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int) {
	client, _ := util.GetHttpClient(infuraKey)
	gasLimit := uint64(21000)
	var data []byte
	signedTx, err := util.SendNewTransaction(client, privateKey, to, value, gasLimit, data)
	if err != nil {
		log.Fatal(err)
	}
	// 打印交易hash
	fmt.Println("tx hash:", signedTx.Hash().Hex())

	receipt, err := util.WaitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("receipt tx hash:", receipt.TxHash.Hex())
	fmt.Println("receipt status:", receipt.Status) // 1为成功，0失败
}

// 部署Counter合约
func deployStoreContract() (*contractbindings.Counter, error) {
	client, _ := util.GetHttpClient(infuraKey)
	privateKey := testPrivateKey
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	address, tx, instance, err := contractbindings.DeployCounter(auth, client)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("contract address:", address.Hex())
	fmt.Println("depoly contract tx hash:", tx.Hash().Hex())

	return instance, nil
}

// 根据合约地址返回实例
func getInstanceByContractAddress(contractAddr common.Address) (*contractbindings.Counter, error) {
	client, _ := util.GetHttpClient(infuraKey)
	counterContract, err := contractbindings.NewCounter(contractAddr, client)
	if err != nil {
		return nil, err
	}
	return counterContract, nil
}

// 合约调用
func callContract(counterContract *contractbindings.Counter, contractAddr common.Address) {
	client, _ := util.GetHttpClient(infuraKey)
	privateKey := testPrivateKey
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 创建交易器
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	// 构建合约调用的消息体（用于估算 gas）
	contractABI, err := contractbindings.CounterMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}
	data, err := contractABI.Pack("addOne") // 打包调用函数及参数
	if err != nil {
		log.Fatalf("打包函数及参数失败: %v", err)
	}
	msg := ethereum.CallMsg{
		From:  opt.From,      // 交易发起地址
		To:    &contractAddr, // 合约地址
		Data:  data,          // 调用数据（函数选择器+参数）
		Value: big.NewInt(0), // 转账金额（ERC20 转账通常为 0）
	}
	estimatedGas, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("估算 gas 失败: %v", err) // 可能需要检查参数或节点
	}
	opt.GasLimit = estimatedGas * 120 / 100 // 增加 20% 缓冲

	// 调用合约 setItem gas通过预估计算
	tx, err := counterContract.AddOne(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("call addOne, tx hash:", tx.Hash().Hex())
	// 等待交易确认
	receipt, err := util.WaitForReceipt(client, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hast receipt status:", receipt.Status)

	callOpt := &bind.CallOpts{Context: context.Background()}
	// 调用合约 items
	valueInContract, err := counterContract.Count(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("call Count result:", valueInContract.String())
}

type AddCount struct {
	NewCount *big.Int
}

// 监听合约事件
func subscribeContractLogs(contractAddr common.Address, wg *sync.WaitGroup) {
	client, _ := util.GetWSSClient(infuraKey)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Subscribe contract event")
	contractAbi, err := contractbindings.CounterMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}

	var topics []string
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog.BlockHash.Hex())
			fmt.Println(vLog.BlockNumber)
			fmt.Println(vLog.TxHash.Hex())
			var event AddCount
			err := contractAbi.UnpackIntoInterface(&event, "AddCount", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("event AddCount:", event.NewCount.Uint64())
			for i := range vLog.Topics {
				topics = append(topics, vLog.Topics[i].Hex())
			}
			fmt.Println("topics[0]=", topics[0])
			if len(topics) > 1 {
				fmt.Println("index topic:", topics[1:])
			}
			wg.Done()
		}
	}
}

func main() {
	// task 1.1
	// queryBlockInfo(big.NewInt(5671744))

	// task 1.2
	// toAddress := common.HexToAddress("0x0405d109770350d2a26bd7874525945106e306cb")
	// value := big.NewInt(1 * 1e15) // 转0.001个ETH in wei (0.001 eth 15个0)
	// transferETH(testPrivateKey, toAddress, value)

	// fmt.Println("------------------task2-------------------")

	// task 2.1部署合约
	// _, err := deployStoreContract()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// task 2.2如果合约已经部署了，需要根据合约地址，拿取合约实例
	// 上一步2.1部署合约得到合约地址 0x41DAE25BBe43A350fd17A43df8d7003dEC242Cf0
	contractAddr := common.HexToAddress("0x41DAE25BBe43A350fd17A43df8d7003dEC242Cf0")
	counterContract, err := getInstanceByContractAddress(contractAddr)
	if err != nil {
		log.Fatal(err)
	}

	// task 2.3.1监听合约事件
	var wg sync.WaitGroup
	wg.Add(1)
	go subscribeContractLogs(contractAddr, &wg)

	// task 2.3.2调用合约
	callContract(counterContract, contractAddr)

	wg.Wait()
}
