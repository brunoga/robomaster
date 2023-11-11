package logger

import (
	"log/slog"
	"os"
	"sync"

	"github.com/lmittmann/tint"
)

type Logger struct {
	*slog.Logger

	m             sync.Mutex
	previousLevel slog.Level
	levelVar      *slog.LevelVar
}

func New(level slog.Level) *Logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(level)

	opts := &tint.Options{
		Level:     levelVar,
		AddSource: true,
	}

	return &Logger{
		Logger:        slog.New(tint.NewHandler(os.Stdout, opts)),
		previousLevel: level,
		levelVar:      levelVar,
	}
}

func (l *Logger) SetLevel(level slog.Level) {
	l.m.Lock()

	l.previousLevel = l.levelVar.Level()
	l.levelVar.Set(level)

	l.m.Unlock()
}

func (l *Logger) ResetLevel() {
	l.m.Lock()

	currentLevel := l.levelVar.Level()

	l.levelVar.Set(l.previousLevel)
	l.previousLevel = currentLevel

	l.m.Unlock()
}
