package dto

import "time"

type OrderPayload struct {
	OrderId     int       `json:"order_id"`
	UserId      string    `json:"user_id"`
	CartId      int       `json:"cart_id"`
	TotalAmount int       `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
