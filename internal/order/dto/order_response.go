package dto

import "time"

type OrderResponse struct {
	Id          int       `json:"id"`
	UserId      string    `json:"user_id"`
	CartId      int       `json:"cart_id"`
	Status      string    `json:"status"`
	TotalAmount int       `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
