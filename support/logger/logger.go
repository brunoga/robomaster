package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/brunoga/groupfilterhandler"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

const (
	LevelTrace = slog.LevelDebug - 1
)

type Logger struct {
	*slog.Logger

	levelVar *slog.LevelVar
}

func New(level slog.Level, allowedGroups ...string) *Logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)

	output := os.Stdout

	opts := &tint.Options{
		Level:   levelVar,
		NoColor: !isatty.IsTerminal(output.Fd()) || runtime.GOOS == "ios",
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			stringValue := attr.Value.String()
			if len(stringValue) > 100 {
				stringValue = stringValue[:100] + "..."
			}
			return slog.Attr{
				Key:   attr.Key,
				Value: slog.StringValue(stringValue),
			}
		},
	}

	return &Logger{
		Logger: slog.New(groupfilterhandler.New(tint.NewHandler(colorable.NewColorable(output),
			opts), allowedGroups...)),
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
	if !l.Enabled(context.Background(), LevelTrace) {
		return func(args ...any) {}
	}

	// Convert args to []slog.Attr (if needed).
	attrs := slog.Group("", args...).Value.Group()

	l.trace(msg, "TRACE START: ", attrs...)
	return func(args ...any) {
		attrs := slog.Group("", args...).Value.Group()
		l.trace(msg, "TRACE END: ", attrs...)
	}
}

func (l *Logger) trace(msg, prefix string, args ...slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), LevelTrace, prefix+msg, pcs[0])
	r.AddAttrs(args...)
	l.Logger.Handler().Handle(context.Background(), r)
}
