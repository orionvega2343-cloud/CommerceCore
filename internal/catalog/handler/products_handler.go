package handler

import (
	"CommerceCore/internal/catalog/domain"
	"CommerceCore/internal/catalog/dto"
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

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func (h *ProductsHandlerImpl) CreateProduct(c *gin.Context) {
	var req dto.ProductRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}

	ctx := c.Request.Context()
	product := toDomainProduct(&req)

	p, err := h.svc.CreateProduct(ctx, product)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	c.JSON(200, toProductResponse(*p))
}

func (h *ProductsHandlerImpl) GetAllProducts(c *gin.Context) {
	var isActivePtr *bool
	if isActive := c.Query("is_active"); isActive != "" {
		b, err := strconv.ParseBool(isActive)
		if err != nil {
			c.JSON(400, dto.ProductResponse{})
			return
		}
		isActivePtr = &b
	}

	parsedLimit := defaultLimit
	if limit := c.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(400, dto.ProductResponse{})
			return
		}
		parsedLimit = l
	}

	parsedOffset := defaultOffset
	if offset := c.Query("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(400, dto.ProductResponse{})
			return
		}
		parsedOffset = o
	}

	ctx := c.Request.Context()
	products, err := h.svc.GetAllProducts(ctx, isActivePtr, parsedLimit, parsedOffset)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	c.JSON(200, toProductListResponse(products))
}

func (h *ProductsHandlerImpl) GetProductById(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	ctx := c.Request.Context()
	product, err := h.svc.GetProductById(ctx, parsedId)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	c.JSON(200, toProductResponse(*product))
}

func (h *ProductsHandlerImpl) UpdateProduct(c *gin.Context) {
	var req dto.ProductRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}

	ctx := c.Request.Context()
	product := toDomainProduct(&req)

	err = h.svc.UpdateProduct(ctx, product)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	c.JSON(200, toProductResponse(*product))
}

func (h *ProductsHandlerImpl) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}

	ctx := c.Request.Context()
	err = h.svc.DeleteProduct(ctx, parsedId)
	if err != nil {
		c.JSON(400, dto.ProductResponse{})
		return
	}
	c.JSON(200, gin.H{})
}
