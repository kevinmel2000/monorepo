package rdbms

import (
	"context"
	"log"

	"github.com/tokopedia/sqlt"
)

//SetMaxConnection set max conn db
func SetMaxConnection(ctx context.Context, maxnumber int) error {
	for _, val := range dbObject.connectedDbs {
		if val == nil {
			continue
		}	
		val.SetMaxOpenConns(maxnumber)
	}
	return nil
}

// Prepare will fatal all failed prepared query
func Prepare(ctx context.Context, db *sqlt.DB, query string) *sqlt.Stmt {
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to prepare query %s. Error: %s", query, err.Error())
	}
	return &stmt
}

// Preparex will fatal all failed preparedx query
func Preparex(ctx context.Context, db *sqlt.DB, query string) *sqlt.Stmtx {
	stmtx, err := db.PreparexContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to preparex query %s. Error: %s", query, err.Error())
	}
	return stmtx
}
