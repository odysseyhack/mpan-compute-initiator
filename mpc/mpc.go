package mpc

import (
	"log"
	"math/big"
)

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

// QueryComputation specifies what should be computed

func StartQueryListener() chan Query {
	qc := make(chan Query)
	go func() {
		for {
			doQuery(<-qc)
		}
	}()
	return qc
}

func doQuery(query Query) {
	// For now, we just print it (looks nice for the dashboard)
	log.Print("Received a query with ID %v:\n", *query.QueryId)
	log.Printf("  %v(%v, %v)\n", query.QueryType, query.Identifier, query.Attribute)
}
