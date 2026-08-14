package postgres

import (
	"CommerceCore/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// HealthCheck - проверяет валидное соединение с БД,
// если соединение не установлено, устанавливает его через PingContext
func HealthCheck(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Ограничиваем время проверки
		//и пробрасываем вниз через контекст
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			slog.Error("failed to ping database", "error", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "failed",
				"error":  "postgres connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func Connect(cfg config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Name, cfg.Database.Password, cfg.Database.SslMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return db, nil
}
