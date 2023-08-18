#services
server:
	@echo "Starting server"
	@go run main.go
#solidity
sol:
	solc --evm-version london --optimize --abi ./contracts/MySmartContract.sol -o build
	solc --evm-version london --optimize --bin ./contracts/MySmartContract.sol -o build
	abigen --abi=./build/MySmartContract.abi --bin=./build/MySmartContract.bin --pkg=api --out=./api/MySmartContract.go