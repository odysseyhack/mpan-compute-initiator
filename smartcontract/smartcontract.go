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
	"github.com/spf13/viper"
)

func WaitForQuery(queryChan chan mpc.Query) {
	conn, err := ethclient.Dial(viper.GetString("eth.url"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// Instantiate the contract and display its name
	gatekeeper, err := NewGateKeeper(common.HexToAddress(viper.GetString("eth.smartcontractAddress")), conn)

	if err != nil {
		log.Fatalf("Failed to instantiate the Greeting contract: %v", err)
	}

	log.Println("Waiting for smart-contract approved queries.")

	var dummyEventChannel = make(chan *GateKeeperNewDummyTrigger)
	var blockNumber uint64 = 1
	// s := []common.Address{}
	opts := &bind.WatchOpts{}
	opts.Start = &blockNumber
	_, err = gatekeeper.WatchNewDummyTrigger(opts, dummyEventChannel)
	if err != nil {
		log.Fatal("Unable to subscribe to Dummy event!", err)
	}

	var realEventChannel = make(chan *GateKeeperNewTrigger)
	// s := []common.Address{}
	_, err = gatekeeper.WatchNewTrigger(opts, realEventChannel)
	if err != nil {
		log.Fatal("Unable to subscribe to Real event!", err)
	}

	for {

		select {
		case event := <-dummyEventChannel:
			log.Println("Received a dummy query request from smart contract.")

			// Now create a query
			var q mpc.Query
			q.Compute = make([]mpc.QueryComputation, 1)
			q.Compute[0] = mpc.QueryComputation{Function: "AVG", AttributeOwner: string(event.ComputeDataOwner), Attribute: string(event.ComputeAttr)}
			q.For = make([]mpc.QuerySelector, 1)
			q.For[0] = mpc.QuerySelector{AttributeOwner: string(event.SelectorDataOwner), Attribute: string(event.SelectorAttr), Operator: string(event.SelectorOperator), AttributeValue: string(event.SelectorAttrValue)}
			q.ID = event.QueryId.String()
			if !event.MinimalResultSize.IsUint64() {
				log.Fatalf("Minimal result size does not fit in uin64: %v", event.MinimalResultSize.String())
			}
			q.MinimalResultSize = event.MinimalResultSize.Uint64()
			log.Printf("q: %v\n", q)
			if err != nil {
				log.Fatalf("Error decoding json: %v", err)
			}
			queryChan <- q

		case event := <-realEventChannel:
			log.Println("Received a real query request from smart contract.")

			// Now create a query
			var q mpc.Query
			q.Compute = make([]mpc.QueryComputation, 1)
			q.Compute[0] = mpc.QueryComputation{Function: "AVG", AttributeOwner: string(event.ComputeDataOwner), Attribute: string(event.ComputeAttr)}
			q.For = make([]mpc.QuerySelector, 1)
			q.For[0] = mpc.QuerySelector{AttributeOwner: string(event.SelectorDataOwner), Attribute: string(event.SelectorAttr), Operator: string(event.SelectorOperator), AttributeValue: string(event.SelectorAttrValue)}
			q.ID = event.QueryId.String()
			if !event.MinimalResultSize.IsUint64() {
				log.Fatalf("Minimal result size does not fit in uin64: %v", event.MinimalResultSize.String())
			}
			q.MinimalResultSize = event.MinimalResultSize.Uint64()
			log.Printf("q: %v\n", q)
			if err != nil {
				log.Fatalf("Error decoding json: %v", err)
			}
			queryChan <- q
		}

	}
}

// StartLogger starts a logger to submit stuff back to the smart contract
func StartLogger() chan mpc.QueryLog {
	ql := make(chan mpc.QueryLog)
	go func() {
		conn, err := ethclient.Dial(viper.GetString("eth.url"))

		if err != nil {
			log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		}

		gatekeeper, err := NewGateKeeper(common.HexToAddress(viper.GetString("eth.smartcontractAddress")), conn)
		if err != nil {
			log.Fatalf("Failed to instantiate the audit_trail contract: %v", err)
		}
		// opts := bind.TransactOpts{From: "", Signer: "", }
		f, err := os.Open(viper.GetString("eth.keyfile"))
		if err != nil {
			log.Fatalf("Failed to open key file: %v", err)
		}
		txOpts, err := bind.NewTransactor(f, viper.GetString("eth.passphrase"))
		if err != nil {
			log.Fatalf("Failed to create transactor: %v", err)
		}
		for {
			logEntry := <-ql
			queryIDint := new(big.Int)
			queryIDint.SetString(logEntry.QueryID, 10)
			switch logEntry.Type {
			case mpc.QueryLogStart:
				tx, err := gatekeeper.JoinQuery(txOpts, queryIDint, []byte(logEntry.LogEntry))
				if err != nil {
					log.Fatalf("Failed to write to audit log: %v", err)
				}
				log.Printf("Submitted join entry to audit_trail (tx: %v).\n", tx.Hash().String())
			case mpc.QueryLogReject:
				tx, err := gatekeeper.RejectQuery(txOpts, queryIDint, []byte(logEntry.LogEntry))
				if err != nil {
					log.Fatalf("Failed to write to audit log: %v", err)
				}
				log.Printf("Submitted reject entry to audit_trail (tx: %v).\n", tx.Hash().String())
			case mpc.QueryLogResult:
				tx, err := gatekeeper.CompleteQuery(txOpts, queryIDint, []byte(logEntry.LogEntry))
				if err != nil {
					log.Fatalf("Failed to write to audit log: %v", err)
				}
				log.Printf("Submitted result entry to audit_trail (tx: %v).\n", tx.Hash().String())
			}
		}
	}()
	return ql
}
