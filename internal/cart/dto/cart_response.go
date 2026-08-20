package dto

import "time"

type CartResponse struct {
	Id        int                `json:"id"`
	UserID    string             `json:"user_id"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Items     []CartItemResponse `json:"items"`
}
