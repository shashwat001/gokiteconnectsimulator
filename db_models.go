package kiteconnectsimulator

import (
	"time"
)

type DbOrder struct {
	Order
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
}

type DbHolding struct {
	Holding
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
}

// func (*DbHolding) AfterCreateTable(query *bun.CreateTableQuery) error {
// 	_, err := query.DB().NewCreateIndex().
// 		Model((*DbHolding)(nil)).
// 		Index("category_id_idx").
// 		Column("category_id").
// 		Exec(context.Background())
// 	return err
// }
