package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ctxKey -типизированный ключ для ctx.Value, не голая строка,
// чтобы исключить конфликты с другими пакетами
type ctxKey string

var requestIDKey ctxKey = "requestID"

// RequestID - генерирует/пробрасывает сквозной id запроса,
// если входящий X-Request-ID уже есть - переиспользуем, если нет, создаем новый с помощью uuid, кладем id:
// в context для пробрасывания вниз по цепочке вызовов и для обработчика внутри сервиса,
// в заголовок ответа чтобы клиент/поддержка видели id в ответе
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		ctx := context.WithValue(c.Request.Context(), requestIDKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-Id", id)
		c.Next()
	}
}

// Хелпер для чтения request_id из контекста, возвращает "",
//если RequestId middleware не отработал - значение по умолчанию

func IDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
