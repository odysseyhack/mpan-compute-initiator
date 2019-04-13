package mpc

import (
	"encoding/json"
	"flag"
	"log"
	"math/big"
	"os"
)

var output *os.File

type QueryType int

const (
	QUERY_TYPE_INFO = 0
	QUERY_TYPE_CALC = 1
)

// Query is the type of query this MPC system supports
type Query struct {
	QueryType  QueryType
	Identifier int
	Attribute  int
	QueryId    *big.Int
}

// Listener loop just does queries sent to its channel
func StartQueryListener() chan Query {
	qc := make(chan Query)
	go func() {
		for {
			doQuery(<-qc)
		}
	}()
	return qc
}

// Send a query to the mpc node
func doQuery(query Query) {
	// For now, we just print it (looks nice for the dashboard)
	log.Print("Received a query with ID %v:\n", *query.QueryId)
	log.Printf("  %v(%v, %v)\n", query.QueryType, query.Identifier, query.Attribute)

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		log.Printf("!!! Error marshaling to string, %v", err)
		return
	}

	log.Printf("Marshalled: %s", jsonQuery)

	output.Write(append(jsonQuery, '\n'))
}

// In init we make a named pipe for comms to the mpc node and prepare it for writing
func init() {
	filename := flag.String("filename", "/tmp/computeInitiatorOutput", "Named pipe for output to mpc node")
	flag.Parse()

	os.Remove(*filename)
	var err error
	output, err = os.OpenFile(*filename, os.O_WRONLY|os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Can not open named pipe %v, %v", *filename, err)
	}
}
