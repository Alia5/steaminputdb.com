package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Alia5/steaminputdb.com/api/ctx"
)

func ParseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warning", "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type colorHandler struct {
	w     io.Writer
	level slog.Leveler
}

func SetupDefault(logLevel string) *slog.Logger {
	level := ParseLevel(logLevel)
	handler := &colorHandler{w: os.Stdout, level: level}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	buf := strings.Builder{}

	buf.WriteString("\033[90m")
	buf.WriteString(r.Time.Format("2006-01-02T15:04:05.000000Z07:00"))
	buf.WriteString("\033[0m ")

	var color string
	switch {
	case r.Level >= slog.LevelError:
		color = "\033[31m" // Red
	case r.Level >= slog.LevelWarn:
		color = "\033[33m" // Yellow
	case r.Level >= slog.LevelInfo:
		color = "\033[32m" // Green
	case r.Level >= slog.LevelDebug:
		color = "\033[34m" // Blue
	default:
		color = "\033[0m"
	}
	buf.WriteString(color)
	fmt.Fprintf(&buf, "%5s", r.Level.String())
	buf.WriteString("\033[0m")

	buf.WriteString(" ")
	buf.WriteString(r.Message)

	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(" ")
		buf.WriteString("\033[90m")
		buf.WriteString(a.Key)
		buf.WriteString("\033[0m=")

		valStr := a.Value.String()
		switch a.Key {
		case string(ctx.KeyStatusCode):
			if strings.HasPrefix(valStr, "2") {
				buf.WriteString("\033[32m") // Green
			} else if strings.HasPrefix(valStr, "3") {
				buf.WriteString("\033[36m") // Cyan
			} else if strings.HasPrefix(valStr, "4") {
				buf.WriteString("\033[33m") // Yellow
			} else if strings.HasPrefix(valStr, "5") {
				buf.WriteString("\033[31m") // Red
			}
		case "err", "error":
			buf.WriteString("\033[31m") // Red for errors
		case "query":
			buf.WriteString("\033[36m") // Cyan for SQL
		default:
			buf.WriteString("\033[37m") // Light gray for other values
		}
		buf.WriteString(valStr)
		buf.WriteString("\033[0m")
		return true
	})

	buf.WriteString("\n")
	_, err := h.w.Write([]byte(buf.String()))
	return err
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *colorHandler) WithGroup(name string) slog.Handler {
	return h
}
