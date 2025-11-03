## 使用 nodejs，安装 solc 工具：
```shell
npm install -g solc
```

### 使用命令，编译合约代码，会在当目录下生成一个编译好的二进制字节码文件 store_sol_Store.bin：
```shell
solcjs --bin ./contracts/Store.sol
```

### 使用命令，生成合约 abi 文件，会在当目录下生成 store_sol_Store.abi 文件：
```shell
solcjs --abi ./contracts/Store.sol
```

## abigin 工具可以使用下面的命令安装：
```shell
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

### 使用 abigen 工具根据这两个生成 bin 文件和 abi 文件，生成 go 代码：
```shell
abigen --bin=contracts_Store_sol_Store.bin --abi=contracts_Store_sol_Store.abi --pkg=store --out=./study/store/store.go
```

# ---------------------------------------TASK--------------------------------------

## 任务1下第二个任务：
### 先把 Counter.sol 生成abi和bin文件
```shell
solcjs --abi --bin ./contracts/Counter.sol
```
生成好后，把这两个abi和bin修改为 Counter.abi 和 Counter.bin

### 可以通过 go generate 代码生成工具，生成go代码文件
```shell
go generate ./internal/contractbindings/...
```
