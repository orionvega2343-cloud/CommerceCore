package dto

type ProductRequest struct {
	Name          string `json:"name"`
	Price         int    `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
	IsActive      bool   `json:"is_active"`
}
