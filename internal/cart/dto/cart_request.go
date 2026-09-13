package dto

type CartRequest struct {
	UserId string            `json:"user_id"`
	Status string            `json:"status"`
	Items  []CartItemRequest `json:"items"`
}
