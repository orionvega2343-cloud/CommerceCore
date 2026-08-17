package dto

type CartItemRequest struct {
	CartId        int `json:"cart_id"`
	ProductId     int `json:"product_id"`
	Quantity      int `json:"quantity"`
	PriceSnapshot int `json:"price_snapshot"`
}
