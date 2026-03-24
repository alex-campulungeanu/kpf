package dlogger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type TeeHandler struct {
	handlers []slog.Handler
}

func (t *TeeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range t.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (t *TeeHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range t.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (t *TeeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &TeeHandler{handlers: newHandlers}
}

func (t *TeeHandler) WithGroup(group string) slog.Handler {
	newHandlers := make([]slog.Handler, len(t.handlers))
	for i, h := range t.handlers {
		newHandlers[i] = h.WithGroup(group)
	}
	return &TeeHandler{handlers: newHandlers}
}

func IsGoRun(path string) bool {
	return strings.Contains(path, os.TempDir())
}

func InitLogger(filePath string, level slog.Level) (*os.File, error) {

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{AddSource: true})
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		AddSource:  true,
	})

	var handlers []slog.Handler
	handlers = append(handlers, consoleHandler)

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if IsGoRun(exe) {
		handlers = append(handlers, fileHandler)
	}

	tee := &TeeHandler{
		handlers: handlers,
	}

	slog.SetDefault(slog.New(tee))
	return file, nil

}
