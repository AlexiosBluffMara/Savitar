package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents the severity of a log entry.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Entry is a single structured log record.
type Entry struct {
	Timestamp string `json:"ts"`
	Level     Level  `json:"level"`
	Component string `json:"component,omitempty"`
	Surface   string `json:"surface,omitempty"`
	SessionID string `json:"session,omitempty"`
	Message   string `json:"msg"`
}

// Logger writes structured JSON log entries to an io.Writer.
type Logger struct {
	mu        sync.Mutex
	out       io.Writer
	component string
	surface   string
	sessionID string
	minLevel  Level
}

var global = &Logger{out: os.Stderr, minLevel: LevelInfo}

// New returns a logger that writes to out at info level or above.
func New(out io.Writer) *Logger {
	return &Logger{out: out, minLevel: LevelInfo}
}

// With returns a child logger with the given fields set.
func (l *Logger) With(component, surface, sessionID string) *Logger {
	return &Logger{
		out:       l.out,
		component: component,
		surface:   surface,
		sessionID: sessionID,
		minLevel:  l.minLevel,
	}
}

// WithComponent returns a child logger with the component field set.
func (l *Logger) WithComponent(component string) *Logger {
	return l.With(component, l.surface, l.sessionID)
}

// SetLevel sets the minimum level. Entries below this level are dropped.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

func (l *Logger) enabled(level Level) bool {
	order := map[Level]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
	}
	return order[level] >= order[l.minLevel]
}

func (l *Logger) write(level Level, msg string) {
	if !l.enabled(level) {
		return
	}
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Component: l.component,
		Surface:   l.surface,
		SessionID: l.sessionID,
		Message:   msg,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintf(l.out, "%s\n", data)
}

func (l *Logger) Debug(msg string) { l.write(LevelDebug, msg) }
func (l *Logger) Info(msg string)  { l.write(LevelInfo, msg) }
func (l *Logger) Warn(msg string)  { l.write(LevelWarn, msg) }
func (l *Logger) Error(msg string) { l.write(LevelError, msg) }

func (l *Logger) Debugf(format string, args ...any) { l.write(LevelDebug, fmt.Sprintf(format, args...)) }
func (l *Logger) Infof(format string, args ...any)  { l.write(LevelInfo, fmt.Sprintf(format, args...)) }
func (l *Logger) Warnf(format string, args ...any)  { l.write(LevelWarn, fmt.Sprintf(format, args...)) }
func (l *Logger) Errorf(format string, args ...any) { l.write(LevelError, fmt.Sprintf(format, args...)) }

// Global logger helpers. These are thin wrappers around the global logger.
func SetOutput(out io.Writer) { global.mu.Lock(); global.out = out; global.mu.Unlock() }
func SetLevel(level Level)    { global.SetLevel(level) }

func Debug(msg string)                   { global.Debug(msg) }
func Info(msg string)                    { global.Info(msg) }
func Warn(msg string)                    { global.Warn(msg) }
func Error(msg string)                   { global.Error(msg) }
func Debugf(format string, args ...any)  { global.Debugf(format, args...) }
func Infof(format string, args ...any)   { global.Infof(format, args...) }
func Warnf(format string, args ...any)   { global.Warnf(format, args...) }
func Errorf(format string, args ...any)  { global.Errorf(format, args...) }

// Writer returns an io.Writer that writes each line as an info-level log entry
// with the given component. This adapts the Logger for use as an io.Writer
// (e.g., as the logger for discord.Bot).
func (l *Logger) Writer() io.Writer {
	return &lineWriter{logger: l}
}

type lineWriter struct {
	logger *Logger
}

func (w *lineWriter) Write(p []byte) (int, error) {
	msg := string(p)
	// Strip trailing newline so the log entry is clean.
	for len(msg) > 0 && (msg[len(msg)-1] == '\n' || msg[len(msg)-1] == '\r') {
		msg = msg[:len(msg)-1]
	}
	if msg != "" {
		w.logger.Info(msg)
	}
	return len(p), nil
}
