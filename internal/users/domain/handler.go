package domain

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserById(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdateRole(c *gin.Context)
}
