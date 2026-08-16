package dto

type ProductResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
	IsActive      bool   `json:"is_active"`
}
