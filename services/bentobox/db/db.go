package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	gorp "gopkg.in/gorp.v1"
)

type Options struct {
	User     string
	Password string
	DBName   string
}

type LastBlock struct {
	NumberId          int64 `db:"number_id"`
	ReceivedTimestamp int64 `db:"received_timestamp"`
}

type Block struct {
	BlockNumberId int64          `db:"block_number_id"`
	BlockNumber   sql.NullString `db:"block_number"`
	BlockHash     sql.NullString `db:"block_hash"`
}

type Block struct {
	BlockNumberId int64          `db:"block_number_id"`
	BlockNumber   sql.NullString `db:"block_number"`
	BlockHash     sql.NullString `db:"block_hash"`
}

type Transaction struct {
	TransactionHash  sql.NullString `db:"transaction_hash"`
	BlockNumber      sql.NullString `db:"tx_block_number"`
	TransactionIndex sql.NullString `db:"transaction_index"`
	From             sql.NullString `db:"tx_from"`
	To               sql.NullString `db:"tx_to"`
}

type Log struct {
	Id              int64          `db:"id"`
	TransactionHash sql.NullString `db:"log_transaction_hash"`
	Data            sql.NullString `db:"data"`
	LogIndex        sql.NullString `db:"log_index"`
	Type            sql.NullString `db:"mined"`
}

type Topic struct {
	LogId   int64          `db:"log_id"`
	Content sql.NullString `db:"content"`
}

func InitDb(options Options) *gorp.DbMap {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		options.User, options.Password, options.DBName)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("sql.Open failed %v", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(Block{}, "blocks")
	dbmap.AddTableWithName(Transaction{}, "transactions")
	dbmap.AddTableWithName(Log{}, "logs").SetKeys(true, "Id")
	dbmap.AddTableWithName(Topic{}, "topics")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalf("Create tables failed %v", err)
	}

	return dbmap
}
