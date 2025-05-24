package logger

import (
	"log/slog"
	"os"
)

type Logger = slog.Logger

func New() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	return logger
}
