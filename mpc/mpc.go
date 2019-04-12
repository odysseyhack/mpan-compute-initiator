package mpc

import (
	"log"
)

// Query is the type of query this MPC system supports
type Query struct {
	Compute           []QueryComputation
	MinimalResultSize uint64
	ID                string
}

// QueryComputation specifies what should be computed
type QueryComputation struct {
	Function   string
	Identifier string
}

// QueryLog contains a log entry for a query
type QueryLog struct {
	QueryID  string
	Type     QueryLogType
	LogEntry string
}

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
	log.Print("Received a query with ID %v:\n", query.ID)
	for _, computation := range query.Compute {
		log.Printf("  %v(%v)\n", computation.Function, computation.Identifier)
	}
}
