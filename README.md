## 1. Initialize project
```
go mod init go-ethereum-sample
```
OR with Github so that it could be used by others

```bash
go mod init github.com/iljapavlovs/go-ethereum-sample
```
This will create blank `go.mod` file for all the dependencies 


## 2. Download GO-Ethereum Dependency  

```bash
go get -d github.com/ethereum/go-ethereum
```
This will add dependencies to ```go.mod``` file , transitive (indirect) as well 

[//]: # (https://geth.ethereum.org/docs/dapp/native-bindings)
## 2. 

### 2.1. Install abigen tool for generating Go binding from Solidity contracts
```bash
$ cd $GOPATH/src/github.com/ethereum/go-ethereum
$ go build ./cmd/abigen
```
### 2.2 Instal Solidity compiler `solc`
https://docs.soliditylang.org/en/v0.8.17/installing-solidity.html#macos-packages

## 3. Generate ABI Specification from the contract using Solidity compiler `solc`
```bash
solc --abi --bin --ast-compact-json --asm contracts/Storage.sol -o build
```
Generates JSON ABI spec `Storage.abi` in `build/` folder
* Storage.bin
* Storage.evm
* Storage.sol_json.ast
## 4. Generate Go bindings (Storage.go) from ABI specs
```bash
abigen --abi build/Storage.abi --pkg main --type Storage --out Storage.go
```