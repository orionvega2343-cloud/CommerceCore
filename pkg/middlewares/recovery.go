package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery - иерархически ставится всегда первым в списке,
// для отслеживания паники нижних обработчиков, сервисов, репозиториев,
// по цепочке логирует то, чем паниковали + стек, отдает клиенту 500
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"panic", rec,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"stack", string(debug.Stack()),
				)
				c.JSON(http.StatusInternalServerError, gin.H{})
			}
		}()
		c.Next()
	}
}
