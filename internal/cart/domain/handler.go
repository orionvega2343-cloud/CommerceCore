package domain

import (
	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	CreateOrGet(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)

	GetGuestCart(c *gin.Context)
	AddItemToCart(c *gin.Context)
	UpdateItemIntoCart(c *gin.Context)
	DeleteItemIntoCart(c *gin.Context)
	MergeGuestCart(c *gin.Context)
}
