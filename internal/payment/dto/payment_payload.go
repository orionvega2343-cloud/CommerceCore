package dto

// PaymentPayload - узкий payload события payment.succeeded. Не отдаём наружу
// сам domain.Payment: у события свой контракт, независимый от внутренней
// модели платежа (тот же принцип, что и OrderPayload у order).
type PaymentPayload struct {
	PaymentId int    `json:"payment_id"`
	OrderId   int    `json:"order_id"`
	Amount    int    `json:"amount"`
	Method    string `json:"method"`
}
