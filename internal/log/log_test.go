package log

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestLoggerWritesJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(buf)
	l.Info("hello world")

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected valid JSON entry: %v (got %q)", err, buf.String())
	}
	if entry.Level != LevelInfo {
		t.Errorf("expected level %q, got %q", LevelInfo, entry.Level)
	}
	if entry.Message != "hello world" {
		t.Errorf("expected message %q, got %q", "hello world", entry.Message)
	}
	if entry.Timestamp == "" {
		t.Error("expected timestamp to be set")
	}
}

func TestLoggerWithComponent(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(buf).WithComponent("discord")
	l.Warn("test")

	var entry Entry
	_ = json.Unmarshal(buf.Bytes(), &entry)
	if entry.Component != "discord" {
		t.Errorf("expected component %q, got %q", "discord", entry.Component)
	}
	if entry.Level != LevelWarn {
		t.Errorf("expected level warn, got %q", entry.Level)
	}
}

func TestLoggerRespectMinLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(buf)
	l.SetLevel(LevelError)
	l.Info("should be dropped")
	l.Debug("also dropped")
	l.Error("this should appear")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d: %q", len(lines), buf.String())
	}
	var entry Entry
	_ = json.Unmarshal([]byte(lines[0]), &entry)
	if entry.Level != LevelError {
		t.Errorf("expected error level, got %q", entry.Level)
	}
}

func TestLineWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(buf)
	w := l.Writer()
	_, _ = w.Write([]byte("line from writer\n"))

	var entry Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if entry.Message != "line from writer" {
		t.Errorf("expected message %q, got %q", "line from writer", entry.Message)
	}
}
