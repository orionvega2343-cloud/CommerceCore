package dto

// PaymentRequest - order_id берётся из URL, не отсюда; status клиент не задаёт -
// он всегда "succeeded", это решает сервис, а не запрос.
type PaymentRequest struct {
	Amount int    `json:"amount"`
	Method string `json:"method"`
}
