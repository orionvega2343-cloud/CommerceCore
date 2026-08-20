package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserIdFromContext - получение userId из контекста
func UserIdFromContext(c *gin.Context) (string, bool) {
	val, ok := c.Get("user_id")
	if !ok {
		return "", false
	}
	userId, ok := val.(string)
	if !ok {
		return "", false
	}
	return userId, true
}

func GuestIdFromCookie(c *gin.Context) string {
	val, err := c.Cookie("guest_id")
	if err != nil || val == "" {
		val = uuid.New().String()
		c.SetCookie("guest_id", val, int((7 * 24 * time.Hour).Seconds()), "/", "", false, true)
	}
	return val
}
