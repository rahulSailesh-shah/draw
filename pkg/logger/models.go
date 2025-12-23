package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Level = zerolog.Level

const (
	LevelDebug Level = Level(zerolog.DebugLevel)
	LevelInfo  Level = Level(zerolog.InfoLevel)
	LevelWarn  Level = Level(zerolog.WarnLevel)
	LevelError Level = Level(zerolog.ErrorLevel)
)

type LogWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Data    map[string]any
}

func toRecord(level Level, msg string, atts ...any) Record {
	data := make(map[string]any)
	for i := 0; i < len(atts); i += 2 {
		if i+1 >= len(atts) {
			break
		}
		key := atts[i].(string)
		value := atts[i+1]
		data[key] = value
	}
	return Record{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
		Data:    data,
	}
}

type EventFn func(ctx context.Context, record Record)
