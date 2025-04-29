package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Log = slog.New(handler)
}