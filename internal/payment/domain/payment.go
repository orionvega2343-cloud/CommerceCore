package domain

import "time"

type Payment struct {
	Id        int       `db:"id"`
	OrderId   int       `db:"order_id"`
	Amount    int       `db:"amount"`
	Status    string    `db:"status"`
	Method    string    `db:"method"`
	CreatedAt time.Time `db:"created_at"`
}
