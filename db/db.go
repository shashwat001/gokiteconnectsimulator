package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/kiteconnectsimulator/models"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DbClient struct {
	Db *bun.DB
}

func (dbClient *DbClient) Connect_db() {

	dbClient.connect_pgx()
	// db.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithVerbose(true),
	// 	bundebug.FromEnv("BUNDEBUG"),
	// ))

}

func (dbClient *DbClient) connect_pgx() {
	config, err := pgx.ParseConfig("postgres://postgres:@localhost:5432/simulator?sslmode=disable")
	if err != nil {
		panic(err)
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	dbClient.Db = bun.NewDB(sqldb, pgdialect.New())
}

func (dbClient *DbClient) connect_inbuilt() {
	dsn := "postgres://postgres:@localhost:5432/simulator?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	dbClient.Db = bun.NewDB(sqldb, pgdialect.New())
}

func (dbClient *DbClient) Create_tables() {
	if dbClient.Db == nil {
		panic("Database not initialized")
	}

	dbClient.Db.ResetModel(context.Background())

	dbClient.Db.NewCreateTable().Model((*DbOrder)(nil)).IfNotExists().Exec(context.Background())
	dbClient.Db.NewCreateTable().Model((*DbHolding)(nil)).IfNotExists().Exec(context.Background())

	fmt.Println("Tables created")
}

func (dbClient *DbClient) Complete_order_and_update_holding(orderPK int64) {
	order := new(DbOrder)
	err := dbClient.Db.NewSelect().Model(order).Where("id = ?", orderPK).Scan(context.Background())

	if err != nil {
		log.Fatal("Error in selecting order from db_order: ", err)
	}

	holding := new(DbHolding)
	count, err := dbClient.Db.NewSelect().
		Model(holding).
		Where("exchange = ?", order.Exchange).
		Where("tradingsymbol = ?", order.TradingSymbol).
		ScanAndCount(context.Background())

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	err = dbClient.Db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewUpdate().Model(order).Set("status = ?", "COMPLETE").Where("id = ?", orderPK).Exec(ctx)

		if err != nil {
			log.Fatal(err)
		}

		if count == 0 {

			holding := &DbHolding{Holding: models.Holding{
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

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

}
