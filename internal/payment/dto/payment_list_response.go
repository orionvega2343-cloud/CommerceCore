package dto

type PaymentListResponse struct {
	Payments []*PaymentResponse `json:"payments"`
}
