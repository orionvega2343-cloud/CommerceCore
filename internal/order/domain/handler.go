package domain

import (
	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	Checkout(c *gin.Context)
	GetOrder(c *gin.Context)
	ListOrders(c *gin.Context)
	TransitionStatus(c *gin.Context)
}
