package domain

import (
	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	Create(c *gin.Context)
	GetPaymentById(c *gin.Context)
	ListPayments(c *gin.Context)
	TotalAmountPayments(c *gin.Context)
}
