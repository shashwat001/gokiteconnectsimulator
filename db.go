package kiteconnectsimulator

import (
	"context"
	"database/sql"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func connect_db() {

	connect_pgx()
	// db.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithVerbose(true),
	// 	bundebug.FromEnv("BUNDEBUG"),
	// ))

}

func connect_pgx() {
	config, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/simulator?sslmode=disable")
	if err != nil {
		panic(err)
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	db = bun.NewDB(sqldb, pgdialect.New())
}

func connect_inbuilt() {
	dsn := "postgres://postgres:@localhost:5432/simulator?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db = bun.NewDB(sqldb, pgdialect.New())
}

func create_tables() {
	if db == nil {
		panic("Database not initialized")
	}

	db.ResetModel(context.Background())

	db.NewCreateTable().Model((*DbOrder)(nil)).IfNotExists().Exec(context.Background())
	db.NewCreateTable().Model((*DbHolding)(nil)).IfNotExists().Exec(context.Background())

	fmt.Println("Tables created")
}

func update_holding(tx bun.Tx, order *DbOrder) {
	holding := new(DbHolding)
	count, err := tx.NewSelect().
		Model(holding).
		Where("exchange = ?", order.Exchange).
		Where("tradingsymbol = ?", order.TradingSymbol).
		ScanAndCount(context.Background())

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if count == 0 {

		holding := &DbHolding{Holding: Holding{
			Tradingsymbol:   order.TradingSymbol,
			Exchange:        order.Exchange,
			InstrumentToken: order.InstrumentToken,
			Quantity:        int(order.Quantity),
			Price:           order.Price,
		}}

		res, err := tx.NewInsert().
			Model(holding).
			Exec(context.Background())

		if err != nil {
			panic(err)
		}

		count, err := res.RowsAffected()

		if err != nil {
			panic(err)
		}

		fmt.Printf("New holdings inserted %d", count)

	}

	if order.TransactionType == "BUY" {

		newTotalQuantity := holding.Quantity + int(order.Quantity)
		newTotalPrice := float64(holding.Quantity)*holding.Price + order.Quantity*order.Price

		holding.Quantity = newTotalQuantity
		holding.Price = newTotalPrice / float64(newTotalQuantity)
	} else if order.TransactionType == "SELL" {

		newTotalQuantity := holding.Quantity - int(order.Quantity)
		newTotalPrice := float64(holding.Quantity)*holding.Price - order.Quantity*order.Price

		holding.Quantity = newTotalQuantity
		holding.Price = newTotalPrice / float64(newTotalQuantity)
	}

	res, err := tx.NewUpdate().
		Model(holding).
		Where("exchange = ?", holding.Exchange).
		Where("tradingsymbol = ?", holding.Tradingsymbol).
		Exec(context.Background())

	if err != nil {
		panic(err)
	}

	updatedRows, err := res.RowsAffected()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Holdings updated %d", updatedRows)
}