package logger

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fatih/color"
)

const runtimeCallerSkip = 4

var currentDir = "."

type Handler struct {
	slog.Handler
	lConsole *log.Logger
	lFile    *log.Logger
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}
	currentDir = dir
}

func NewHandler(debug bool) *Handler {
	return &Handler{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		lConsole: log.New(os.Stdout, "", 0),
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	err := h.console(r)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) console(r slog.Record) error {
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05]")
	_, file, line, ok := runtime.Caller(runtimeCallerSkip)
	if ok {
		relPath, err := filepath.Rel(currentDir, file)
		if err != nil {
			relPath = file
		}
		h.lConsole.Printf("%s %s %s:%s [ %s ] %s \n",
			color.GreenString(timeStr),
			level,
			color.HiCyanString(relPath),
			color.HiCyanString(strconv.Itoa(line)),
			color.CyanString(r.Message),
			color.HiWhiteString(string(b)))
	} else {
		h.lConsole.Println(color.GreenString(timeStr),
			level,
			color.CyanString(r.Message),
			color.WhiteString(string(b)))
	}

	return nil
}

type Logger struct {
	*slog.Logger
}

func NewLogger(debug bool) *Logger {
	handler := NewHandler(debug)
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}
