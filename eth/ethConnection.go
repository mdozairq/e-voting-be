package eth

import "github.com/ethereum/go-ethereum/ethclient"

func GetEthClient() *ethclient.Client {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		panic(err)
	}
	return client
}
