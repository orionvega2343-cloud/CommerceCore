package handler

import (
	"CommerceCore/internal/catalog/domain"
	"CommerceCore/internal/catalog/domain/errs"
	"CommerceCore/internal/catalog/dto"
	"CommerceCore/pkg/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ domain.ProductHandler = (*ProductsHandlerImpl)(nil)

type ProductsHandlerImpl struct {
	svc domain.ProductService
}

func NewProductsHandlerImpl(svc domain.ProductService) *ProductsHandlerImpl {
	return &ProductsHandlerImpl{svc: svc}
}

func toDomainProduct(req *dto.ProductRequest) *domain.Product {
	return &domain.Product{
		Name:          req.Name,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		IsActive:      req.IsActive,
	}
}

func toProductResponse(p domain.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		Id:            p.Id,
		Name:          p.Name,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		IsActive:      p.IsActive,
	}
}

func toProductListResponse(products []*domain.Product) *dto.ProductListResponse {
	res := make([]*dto.ProductResponse, 0, len(products))
	for _, p := range products {
		res = append(res, toProductResponse(*p))
	}
	return &dto.ProductListResponse{Products: res}
}

// mapServiceError - переводит доменную ошибку сервиса в HTTP-статус
func mapServiceError(err error) int {
	switch {
	case errors.Is(err, errs.InvalidProductName), errors.Is(err, errs.InvalidPrice), errors.Is(err, errs.InvalidQuantity):
		return http.StatusBadRequest
	case errors.Is(err, errs.InvalidRole):
		return http.StatusForbidden
	case errors.Is(err, errs.ProductNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func (h *ProductsHandlerImpl) CreateProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: "failed to bind request", Code: "FAILED_TO_BIND"})
		return
	}

	ctx := c.Request.Context()
	product := toDomainProduct(&req)
	role := c.GetString("role")

	p, err := h.svc.CreateProduct(ctx, product, role)
	if err != nil {
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_CREATE_PRODUCT"})
		return
	}
	c.JSON(200, toProductResponse(*p))
}

func (h *ProductsHandlerImpl) GetAllProducts(c *gin.Context) {
	var isActivePtr *bool
	if isActive := c.Query("is_active"); isActive != "" {
		b, err := strconv.ParseBool(isActive)
		if err != nil {
			c.JSON(400, response.Error{Message: "failed to parse is_active", Code: "FAILED_TO_GET_PRODUCT"})
			return
		}
		isActivePtr = &b
	}

	parsedLimit := defaultLimit
	if limit := c.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(400, response.Error{Message: "failed to parse limit", Code: "FAILED_TO_PARSE_LIMIT"})
			return
		}
		parsedLimit = l
	}

	parsedOffset := defaultOffset
	if offset := c.Query("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(400, response.Error{Message: "failed to parse offset", Code: "FAILED_TO_PARSE_OFFSET"})
			return
		}
		parsedOffset = o
	}

	ctx := c.Request.Context()
	products, err := h.svc.GetAllProducts(ctx, isActivePtr, parsedLimit, parsedOffset)
	if err != nil {
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_GET_PRODUCTS"})
		return
	}
	c.JSON(200, toProductListResponse(products))
}

func (h *ProductsHandlerImpl) GetProductById(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, response.Error{Message: "failed to parse id", Code: "FAILED_TO_PARSE_ID"})
		return
	}
	ctx := c.Request.Context()
	product, err := h.svc.GetProductById(ctx, parsedId)
	if err != nil {
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_GET_PRODUCT_BY_ID"})
		return
	}
	c.JSON(200, toProductResponse(*product))
}

func (h *ProductsHandlerImpl) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, response.Error{Message: "failed to parse id", Code: "FAILED_TO_PARSE_ID"})
		return
	}

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: "failed to bind request", Code: "FAILED_TO_BIND"})
		return
	}

	ctx := c.Request.Context()
	product := toDomainProduct(&req)
	product.Id = parsedId

	err = h.svc.UpdateProduct(ctx, product)
	if err != nil {
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_UPDATE_PRODUCT"})
		return
	}
	c.JSON(200, toProductResponse(*product))
}

func (h *ProductsHandlerImpl) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, response.Error{Message: "failed to parse id", Code: "FAILED_TO_PARSE_ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.svc.DeleteProduct(ctx, parsedId)
	if err != nil {
		c.JSON(mapServiceError(err), response.Error{Message: err.Error(), Code: "FAILED_TO_DELETE_PRODUCT"})
		return
	}
	c.JSON(200, gin.H{})
}
