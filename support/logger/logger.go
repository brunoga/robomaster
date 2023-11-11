package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
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
		NoColor: !isatty.IsTerminal(output.Fd()),
	}

	return &Logger{
		Logger: slog.New(tint.NewHandler(colorable.NewColorable(output),
			opts)),
		levelVar: levelVar,
	}
}
