package logger

import (
	"context"
	"io"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/rs/zerolog"
)

type Config struct {
	MinLevel  Level
	Service   string
	Handlers  []io.WriteCloser
	TraceIDFn func(ctx context.Context) string
}

type Logger struct {
	zlog          zerolog.Logger
	eventHandlers map[Level]EventFn
	handlers      []io.WriteCloser
	traceIDFn     func(ctx context.Context) string
}

func NewLogger(cfg Config) *Logger {
	writers := make([]io.Writer, len(cfg.Handlers))
	for i, writer := range cfg.Handlers {
		writers[i] = writer
	}

	zlogger := zerolog.New(io.MultiWriter(writers...)).
		Level(cfg.MinLevel).
		With().
		Timestamp().
		Str("service", cfg.Service).
		Logger()

	return &Logger{
		zlog:          zlogger,
		eventHandlers: make(map[Level]EventFn),
		handlers:      cfg.Handlers,
		traceIDFn:     cfg.TraceIDFn,
	}
}

func (l *Logger) On(level Level, fn EventFn) {
	l.eventHandlers[level] = fn
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelInfo, 3, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelError, 3, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelWarn, 3, msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelDebug, 3, msg, args...)
}

func (l *Logger) log(ctx context.Context, level Level, caller int, message string, args ...any) {
	// Handle Event
	handler, ok := l.eventHandlers[level]
	if ok {
		record := toRecord(level, message, args)
		handler(ctx, record)
	}

	// Get caller
	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])
	file, line := runtime.FuncForPC(pcs[0]).FileLine(pcs[0])

	// Build Log Record
	event := l.zlog.WithLevel(level)
	event.Str("msg", message)
	event.Str("file", filepath.Base(file)+":"+strconv.Itoa(line))
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		key := args[i].(string)
		value := args[i+1]
		event = event.Interface(key, value)
	}
	if l.traceIDFn != nil {
		event.Str("trace_id", l.traceIDFn(ctx))
	}
	// Send Log Record
	event.Send()
}

func (l *Logger) Close() error {
	for _, handler := range l.handlers {
		if err := handler.Close(); err != nil {
			return err
		}
	}
	return nil
}
