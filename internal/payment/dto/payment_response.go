package dto

import "time"

type PaymentResponse struct {
	Id        int       `json:"id"`
	OrderId   int       `json:"order_id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
	CreatedAt time.Time `db:"created_at"`
}
