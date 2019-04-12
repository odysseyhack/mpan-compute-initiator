package smartcontract

//go:generate abigen -abi gate_keeper.abi -pkg smartcontract -type GateKeeper -out gatekeeper.go

import (
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/odysseyhack/mpan-compute-initiator/mpc"
)

const (
	SMARTCONTRACT_ADDRESS = "0x368f79382cc5a7b769134369a2de7f5b97b28041"
	ETHEREUM_URL          = "ws://localhost:8546"
)

func WaitForQuery(queryChan chan mpc.Query) {
	conn, err := ethclient.Dial(viper.GetString("eth.url"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// Instantiate the contract and display its name
	gatekeeper, err := NewGateKeeper(common.HexToAddress(SMARTCONTRACT_ADDRESS), conn)

	if err != nil {
		log.Fatalf("Failed to instantiate the Greeting contract: %v", err)
	}

	log.Println("Waiting for smart-contract approved queries.")

	var realEventChannel = make(chan *GateKeeperNewTrigger)
	var blockNumber uint64 = 1
	opts := &bind.WatchOpts{}
	opts.Start = &blockNumber
	_, err = gatekeeper.WatchNewTrigger(opts, realEventChannel)
	if err != nil {
		log.Fatal("Unable to subscribe to event!", err)
	}

	for {

		event := <-realEventChannel
		log.Println("Received a real query request from smart contract.")

		// Now create a query
		var q mpc.Query
		q.Compute = make([]mpc.QueryComputation, 1)
		q.Compute[0] = mpc.QueryComputation{Function: "Info", Identifier: string(event.Identifier)}
		q.ID = event.QueryId.String()
		if !event.MinimalResultSize.IsUint64() {
			log.Fatalf("Minimal result size does not fit in uin64: %v", event.MinimalResultSize.String())
		}
		q.MinimalResultSize = event.MinimalResultSize.Uint64()
		log.Printf("q: %v\n", q)

		queryChan <- q

	}
}
