package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
}

type logger struct {
	Slog *slog.Logger
}

// New returns a new structured logger.
func New(lvl string) Logger {
	logLevel := &slog.LevelVar{}

	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	if lvl == "debug" {
		logLevel.Set(slog.LevelDebug)
	}

	slog := slog.New(opts.NewJSONHandler(os.Stdout))

	return &logger{Slog: slog}
}

func (l *logger) Debug(msg string, args ...any) {
	if args == nil {
		l.Slog.Debug(msg)
		return
	}

	l.Slog.Debug(msg, slog.Any("args", args))
}

func (l *logger) Info(msg string, args ...any) {
	if args == nil {
		l.Slog.Info(msg)
		return
	}

	l.Slog.Info(msg, slog.Any("args", args))
}

func (l *logger) Error(msg string, args ...any) {
	if args == nil {
		l.Slog.Error(msg)
		return
	}

	l.Slog.Error(msg, slog.Any("args", args))
}

func (l *logger) Fatal(msg string, args ...any) {
	if args == nil {
		l.Slog.Error(msg)

		os.Exit(1)
		return
	}

	l.Slog.Error(msg, slog.Any("args", args))
	os.Exit(1)
}
