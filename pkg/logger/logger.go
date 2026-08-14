package logger

import (
	"log/slog"
	"os"
)

//NewLogger - структурированный JSON логгер,
//для отслеживания событий

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
