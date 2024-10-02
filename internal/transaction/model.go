package transaction

import "time"

type OrderSide int

const (
	Buy OrderSide = iota
	Sell
)

type Transaction struct {
	Id           int64     `db:"id" json:"id"`
	Side         OrderSide `db:"side" json:"side"`
	Quantity     float64   `db:"quantity" json:"quantity"`
	Price        float64   `db:"price" json:"price"`
	AllocationId int64     `db:"allocation_id" json:"allocation_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type AmountAndPLResult struct {
    CurrentAmount float64
    UnrealizedPL  float64
}
