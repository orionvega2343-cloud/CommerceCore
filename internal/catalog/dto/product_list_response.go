package dto

type ProductListResponse struct {
	Products []*ProductResponse `json:"products"`
}
