package main

import (
	"CommerceCore/pkg/config"
	"CommerceCore/pkg/logger"
	"CommerceCore/pkg/middlewares"
	"CommerceCore/pkg/postgres"
	"log"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(logger.NewLogger())
	r := gin.Default()
	cfg := config.MustLoad()
	db, err := postgres.Connect(*cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			slog.Error("failed to close database connection:", "error", cerr)
		}
	}()

	r.Use(middlewares.Recovery())
	r.Use(middlewares.RequestID())
	r.Use(middlewares.Logger())
	r.Use(middlewares.Auth())
	r.GET("/connection", postgres.HealthCheck(db))

	if err := r.Run(":" + strconv.Itoa(cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
