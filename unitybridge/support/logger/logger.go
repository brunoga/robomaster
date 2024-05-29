package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

const (
	LevelTrace = slog.Level(slog.LevelDebug - 1)
)

type Logger struct {
	*slog.Logger

	levelVar *slog.LevelVar
}

func New(level slog.Level) *Logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)

	output := os.Stdout

	opts := &tint.Options{
		Level:   levelVar,
		NoColor: !isatty.IsTerminal(output.Fd()) || runtime.GOOS == "ios",
	}

	return &Logger{
		Logger: slog.New(tint.NewHandler(colorable.NewColorable(output),
			opts)),
		levelVar: levelVar,
	}
}

func (l *Logger) Level() slog.Level {
	return l.levelVar.Level()
}

func (l *Logger) WithGroup(group string) *Logger {
	return &Logger{
		Logger:   l.Logger.WithGroup(group),
		levelVar: l.levelVar,
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger:   l.Logger.With(args...),
		levelVar: l.levelVar,
	}
}

func (l *Logger) Trace(msg string, args ...any) func(args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, "TRACE START: "+msg, args...)
	return func(args ...any) {
		l.Logger.Log(context.Background(), LevelTrace, "TRACE END: "+msg, args...)
	}
}
