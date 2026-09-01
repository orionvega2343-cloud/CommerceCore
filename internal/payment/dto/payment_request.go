package dto

type PaymentRequest struct {
	OderId int    `json:"oder_id"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
	Method string `json:"method"`
}
