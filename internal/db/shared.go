package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Cleanup() {
	pool.Close()
}

func Setup(databaseUrl string) {
	var err error
	pool, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
}

func StartTransaction() (pgx.Tx, error) {
	return pool.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
}

func StartReadOnlyTransaction() (pgx.Tx, error) {
	return pool.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadOnly,
		DeferrableMode: pgx.NotDeferrable,
	})
}

func FinishTransaction(tx pgx.Tx, err error) {
	if err == nil {
		tx.Commit(context.Background())
	} else {
		tx.Rollback(context.Background())
	}
}
