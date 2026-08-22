package dto

type OrderItemResponse struct {
	Id           int `json:"id"`
	OrderId      int `json:"order_id"`
	ProductId    int `json:"product_id"`
	Quantity     int `json:"quantity"`
	PricePerUnit int `json:"price_per_unit"`
}
