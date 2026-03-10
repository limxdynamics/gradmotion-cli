package log

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	writer io.Writer
	mu     sync.Mutex
}

func New(logFile string) (*Logger, io.Closer, error) {
	if logFile == "" {
		return &Logger{writer: os.Stderr}, nil, nil
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}
	w := io.MultiWriter(os.Stderr, f)
	return &Logger{writer: w}, f, nil
}

func (l *Logger) Info(fields map[string]any) {
	l.emit("info", fields)
}

func (l *Logger) Error(fields map[string]any) {
	l.emit("error", fields)
}

func (l *Logger) emit(level string, fields map[string]any) {
	if l == nil || l.writer == nil {
		return
	}
	if fields == nil {
		fields = map[string]any{}
	}
	fields["level"] = level
	fields["time"] = time.Now().Format(time.RFC3339)

	b, err := json.Marshal(fields)
	if err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = l.writer.Write(append(b, '\n'))
}
