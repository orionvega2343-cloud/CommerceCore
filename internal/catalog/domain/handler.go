package domain

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetAllProducts(c *gin.Context)
	GetProductById(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}
