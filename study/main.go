package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/insharecgc/goEthereumTask/internal/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"

	token "github.com/insharecgc/goEthereumTask/study/erc20"
	"github.com/insharecgc/goEthereumTask/study/erc721"
)

const infuraKey = "d56339e6f6a0412ea3ae3710fe72f198"

var testPrivateKey, _ = crypto.HexToECDSA("your secret key")

func main() {
	// queryBlock()
	// queryTx()
	// queryReceipt()
	// createCrypto()
	// transferETH()
	// utilTransferEth()
	// transferERC20()
	// transferERC721()
	// queryERC20Info()
	subscribeHead()
}

// 查询区块信息
func queryBlock() {
	client, _ := util.GetHttpClient(infuraKey)
	blockNumber := big.NewInt(5671744)

	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number)              // 5671744
	fmt.Println(header.Hash())              // 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	fmt.Println(header.Time)                // 1712798400
	fmt.Println(header.Difficulty.Uint64()) // 0

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.Time())                // 1712798400
	fmt.Println(block.Difficulty().Uint64()) // 0
	fmt.Println(block.Hash().Hex())          // 0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5
	fmt.Println(block.Transactions().Len())  // 70
}

// 查询交易信息
func queryTx() {
	client, _ := util.GetHttpClient(infuraKey)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())        // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		fmt.Println(tx.Value().String())    // 100000000000000000
		fmt.Println(tx.Gas())               // 21000
		fmt.Println(tx.GasPrice().Uint64()) // 100000000000
		fmt.Println(tx.Nonce())             // 245132
		fmt.Println(tx.Data())              // []
		fmt.Println(tx.To().Hex())          // 0x8F9aFd209339088Ced7Bc0f57Fe08566ADda3587

		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Println("sender", sender.Hex()) // 0x2CdA41645F2dBffB852a605E92B185501801FC28
		} else {
			log.Fatal(err)
		}

		// 拿到交易hash，查询交易的单据信息
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(receipt.Status)                // 1
		fmt.Println(receipt.Logs)                  // []
		fmt.Println(receipt.TxHash.Hex())          // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
		fmt.Println(receipt.TransactionIndex)      // 0
		fmt.Println(receipt.ContractAddress.Hex()) // 0x0000000000000000000000000000000000000000
		break
	}

	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(isPending)
	fmt.Println(tx.Hash().Hex())
}

// 查询收据信息
func queryReceipt() {
	client, _ := util.GetHttpClient(infuraKey)
	// 根据区域信息，查询收据信息（区块下每一个交易的收据信息）
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	recepitByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		log.Fatal(err)
	}
	for _, receipt := range recepitByHash {
		fmt.Println(receipt.Status)
		fmt.Println(receipt.Logs)
		fmt.Println(receipt.TransactionIndex)
		fmt.Println(receipt.BlockNumber)
		fmt.Println(receipt.GasUsed)
		fmt.Println(receipt.ContractAddress.Hex())
		break
	}

	// 根据交易has，查询收据信息
	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(receipt.Status)                // 1
	fmt.Println(receipt.Logs)                  // []
	fmt.Println(receipt.TxHash.Hex())          // 0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5
	fmt.Println(receipt.TransactionIndex)      // 0
	fmt.Println(receipt.ContractAddress.Hex()) // 0x0000000000000000000000000000000000000000
}

// 创建钱包
func createCrypto() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("privateKey:", privateKey)
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("privateKeyBytes:", hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'

	// 通过私钥实例拿公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publicKeyBytes:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'
	// 拿公钥地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address:", address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("hash:", hexutil.Encode(hash.Sum(nil)[12:]))
}

func queryRevertReason(txHash common.Hash) {
	client, _ := util.GetHttpClient(infuraKey)
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == types.ReceiptStatusFailed {
		reason, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &receipt.ContractAddress,
			Data: receipt.TxHash.Bytes(),
		}, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("revert reason:", hexutil.Encode(reason))
	}
}

// 调用util封装的交易函数
func utilTransferEth() {
	client, _ := util.GetHttpClient(infuraKey)
	privateKey := testPrivateKey
	toAddress := common.HexToAddress("0x0405d109770350d2a26bd7874525945106e306cb")
	value := big.NewInt(1 * 1e15) // 转0.001个ETH in wei (0.001 eth 15个0)
	gasLimit := uint64(40000)     // in units
	var data []byte
	signedTx, err := util.SendNewTransaction(client, privateKey, toAddress, value, gasLimit, data)
	if err != nil {
		log.Fatal(err)
	}
	// 打印交易hash
	fmt.Println("tx hash: ", signedTx.Hash().Hex())
}

// ETH 转账
func transferETH() {
	client, _ := util.GetHttpClient(infuraKey)
	privateKey := testPrivateKey
	toAddress := common.HexToAddress("0x0405d109770350d2a26bd7874525945106e306cb")
	value := big.NewInt(1 * 1e16) // in wei (0.01 eth 16个0)
	gasLimit := uint64(21000)     // in units
	var data []byte
	// 1. 构造交易
	// 私钥
	// 私钥拿公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress:", fromAddress)
	// 查询账户余额
	fromBalance, err := client.BalanceAt(context.Background(), fromAddress, nil) // nil 表示最新的一个区块
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("fromBalance:", fromBalance)
	fbalance := new(big.Float)
	fbalance.SetString(fromBalance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1e18))
	fmt.Println("fromBalance(eth):", ethValue)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("to:", toAddress, " chainId:", chainId, " nonce:", nonce)
	// 2. 签名交易
	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 验证签名
	recoveredAddr, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Fatal("签名无效:", err)
	}
	if recoveredAddr != fromAddress {
		log.Fatalf("签名地址不匹配：预期%s，实际%s", fromAddress.Hex(), recoveredAddr.Hex())
	}
	// 3. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	// 打印交易hash
	fmt.Println("tx hash: ", signedTx.Hash().Hex())
	// 8. 等待交易回执
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		log.Fatal("等待交易确认失败:", err)
	}
	if receipt.Status == 1 {
		fmt.Println("交易成功！区块号：", receipt.BlockNumber.Uint64())
	} else {
		fmt.Println("交易失败，gasUsed：", receipt.GasUsed)
	}
}

// 构造 ERC20 transfer 方法的调用数据
func ERC20Transfer(to common.Address, amount *big.Int) []byte {
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	selectorId := hash.Sum(nil)[:4]
	fmt.Printf("selectorId:%x\n", selectorId) // 0xa9059cbb
	// 拼接参数：to 地址（32字节） + amount（32字节）
	paddedTo := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data := append(selectorId, paddedTo...)
	data = append(data, paddedAmount...)
	fmt.Printf("data:%x\n", data)
	return data
}

// transferERC20 函数实现ERC20代币转账功能
func transferERC20() {
	// 创建以太坊客户端
	client, _ := util.GetHttpClient(infuraKey)
	tokenAddress := common.HexToAddress("0x1890491cd06bB4de1b74286AC0b704C1241E9c63") // ERC20代币合约地址

	// 私钥
	privateKey := testPrivateKey
	// 私钥拿公钥地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("fromAddress:", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 转账金额和接收方
	value := big.NewInt(0)     // 转代币，eth金额设置0
	gasLimit := uint64(100000) // ERC20 转账建议设置较高的 Gas Limit（可通过 EstimateGas 精确估算）
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0xcc0089b3882bfff3f476d506160c580cf28d9242")

	// 构造交易
	txAmount := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)) // 转账10个代币
	data := ERC20Transfer(toAddress, txAmount)                      // ERC20转账方法数据
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 签名交易
	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 验证签名
	recoveredAddr, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Fatal("签名无效:", err)
	}
	if recoveredAddr != fromAddress {
		log.Fatalf("签名地址不匹配：预期%s，实际%s", fromAddress.Hex(), recoveredAddr.Hex())
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	// 打印交易hash
	fmt.Println("tx hash: ", signedTx.Hash().Hex())

}

// 辅助函数：手动构造 safeTransferFrom 的调用数据（用于 gas 估算）
func erc721EncodeSafeTransferFrom(from, to common.Address, tokenId *big.Int) []byte {
	// safeTransferFrom 方法 selector: 0x42842e0e
	selector := common.Hex2Bytes("42842e0e")
	// 拼接参数：from（32字节） + to（32字节） + tokenId（32字节）
	paddedFrom := common.LeftPadBytes(from.Bytes(), 32)
	paddedTo := common.LeftPadBytes(to.Bytes(), 32)
	paddedTokenId := common.LeftPadBytes(tokenId.Bytes(), 32)
	data := append(selector, paddedFrom...)
	data = append(data, paddedTo...)
	data = append(data, paddedTokenId...)
	fmt.Printf("data:%x\n", data)
	return data
}

// transferERC721 函数实现ERC721代币转账功能
func transferERC721() {
	// 创建以太坊客户端
	client, _ := util.GetHttpClient(infuraKey)
	tokenAddress := common.HexToAddress("0x4E8Ef74A824d4ef1C83D7c231c4bed5f4a0a6115") // ERC721代币合约地址

	// 私钥
	privateKey := testPrivateKey
	// 私钥拿公钥地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("fromAddress:", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0xcc0089b3882bfff3f476d506160c580cf28d9242")
	tokenId := big.NewInt(0)
	// 估算gasLimit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: erc721EncodeSafeTransferFrom(fromAddress, toAddress, tokenId),
	})
	if err != nil {
		log.Printf("估算 gas 失败，使用默认值: %v", err)
		gasLimit = 100000 // 手动设置默认值
	}

	// 签名交易
	chainId, err := client.NetworkID(context.Background())
	// 构造交易
	opts := &bind.TransactOpts{
		From:     fromAddress,
		Nonce:    big.NewInt(int64(nonce)),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			// 用私钥签名交易（需匹配链 ID）
			return types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
		},
	}
	// 初始化nft实例合约
	nft, err := erc721.NewErc721(tokenAddress, client)
	// 调用合约方法转nft
	tx, err := nft.TransferFrom(opts, fromAddress, toAddress, tokenId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NFT交易已提交, tx hash: ", tx.Hash().Hex())

}

// 查询ERC20代币信息
func queryERC20Info() {
	// 创建以太坊客户端
	client, _ := util.GetHttpClient(infuraKey)
	tokenAddress := common.HexToAddress("0x1890491cd06bB4de1b74286AC0b704C1241E9c63") // ERC20代币合约地址

	// 获取合约实例
	contract, err := token.NewErc20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0xcc0089b3882bfff3f476d506160c580cf28d9242")
	balance, err := contract.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	name, err := contract.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := contract.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := contract.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	totalSupply, err := contract.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("name:", name, "symbol:", symbol)
	fmt.Printf("decimals:%v\n", decimals)
	fmt.Printf("balance wei:%s\n", balance)
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	value := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance:%f\n", value)

	fmt.Printf("totalSupply wei:%s\n", totalSupply)
	fTotalSupply := new(big.Float)
	fTotalSupply.SetString(totalSupply.String())
	value = new(big.Float).Quo(fTotalSupply, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("totalSupply:%f\n", value)

}

// 订阅区块
func subscribeHead() {
	client, _ := util.GetWSSClient(infuraKey)
	// 订阅区块
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println("new block hash:", header.Hash().Hex())
			fmt.Println("new block number:", header.Number.Uint64())
			fmt.Println("new block time:", header.Time)
			fmt.Println("new block nonce:", header.Nonce)
			// 获取区块信息
			block, err := client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block number:", block.Number().Uint64())
			fmt.Println("block hash:", block.Hash().Hex())
			fmt.Println("block time:", block.Time())
			fmt.Println("block nonce:", block.Nonce())
			fmt.Println("block len:", block.Transactions().Len())
			fmt.Println("block size:", block.Size())
			fmt.Println("block gas limit:", block.GasLimit())
			fmt.Println("block gas used:", block.GasUsed())
		}
	}
}
