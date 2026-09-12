package services

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type LoggingService struct {
	logger *slog.Logger
}

func NewLoggingService() *LoggingService {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	return &LoggingService{logger: slog.New(handler)}
}

func (s *LoggingService) Error(r *http.Request, msg string, err error) {
	s.logRequest(r, slog.LevelError, msg, err)
}

func (s *LoggingService) logRequest(r *http.Request, level slog.Level, msg string, err error) {
	var attrs []slog.Attr
	ctx := context.Background()
	if r != nil {
		ctx = r.Context()
		attrs = append(attrs,
			slog.String("METHOD", r.Method),
			slog.String("PATH", r.URL.Path),
		)
	}
	if err != nil {
		attrs = append(attrs, slog.Any("ERROR", err))
	}
	s.logger.LogAttrs(ctx, level, msg, attrs...)
}
