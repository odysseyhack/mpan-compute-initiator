package smartcontract

//go:generate abigen -abi governance.abi -pkg smartcontract -type Governance -out governance.go

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event" //temp
	"github.com/odysseyhack/mpan-compute-initiator/mpc"
)

const (
	SMARTCONTRACT_ADDRESS = "0xa652605f3794d5cd868aa5f295e60fae924fe836"
	ETHEREUM_URL          = "ws://127.0.0.1:8546"
)

func WaitForQuery(queryChan chan mpc.Query) {
	conn, err := ethclient.Dial(ETHEREUM_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	governance, err := NewGovernance(common.HexToAddress(SMARTCONTRACT_ADDRESS), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate the Governance contract: %v", err)
	}

	computeAddress, err := governance.GetGatekeeperAddress(nil)
	if err != nil {
		log.Fatalf("Failed to talk to the Governance contract: %v", err)
	}

	// Instantiate the contract and display its name
	computeContract, err := NewComputeContract(computeAddress, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate the Compute contract: %v", err)
	}

	log.Println("Waiting for smart-contract approved queries.")

	var realEventChannel = make(chan *ComputeContractNewTrigger)
	var blockNumber uint64 = 1
	opts := &bind.WatchOpts{}
	opts.Start = &blockNumber
	_, err = computeContract.WatchNewTrigger(opts, realEventChannel)
	if err != nil {
		log.Fatal("Unable to subscribe to event!", err)
	}

	for {

		event := <-realEventChannel
		log.Println("(SmartContract) Received a real query request from smart contract.")

		// Now create a query
		q := mpc.Query{
			QueryType:       mpc.QueryType(event.QueryType),
			Identifier:      int(event.Identifier),
			Attribute:       int(event.Attribute),
			ClientReference: event.ClientReference,
		}

		log.Printf("(SmartContract) q: %v\n", q)

		queryChan <- q

	}
}

// temporary
type ComputeContract struct{}
type ComputeContractNewTrigger struct {
	QueryType       int
	Identifier      int
	Attribute       int
	ClientReference string
}

func NewComputeContract(address common.Address, backend bind.ContractBackend) (*ComputeContract, error) {
	return nil, nil
}

func (_ComputeContract *ComputeContract) WatchNewTrigger(opts *bind.WatchOpts, sink chan<- *ComputeContractNewTrigger) (event.Subscription, error) {
	return nil, nil
}
