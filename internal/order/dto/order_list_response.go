package dto

type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
}
