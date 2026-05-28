package config

import (
	"log/slog"
	"os"
)

func NewLogger(level string) *slog.Logger {
	var l slog.Level
    
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	default:
		l = slog.LevelInfo
	}
	
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}