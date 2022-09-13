package db

import (
	"kiteconnectsimulator/models"
	"time"
)

type DbOrder struct {
	models.Order
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
}

type DbHolding struct {
	models.Holding
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,default:now()"`
}
