package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims - структура полезной нагрузки JWT токена,
// содержит user_id и стандартные поля (exp, iat и тп)
type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// Auth - middleware для проверки JWT токена,
// ожидает заголовок Authorization: Bearer <token>,
// при успешной валидации кладёт user_id в контекст запроса,
// при ошибке — прерывает цепочку и возвращает 401
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "No token provided"})
			return
		}

		//Проверяем формат Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		tokenString := parts[1]

		//Парсим и валидируем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		//Обработка прочих ошибок или не валидного токена
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
