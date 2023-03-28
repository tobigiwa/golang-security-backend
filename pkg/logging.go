package logging

import (
	"encoding/json"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type LogLevel int8

const (
	LevelInfo LogLevel = iota
	LevelError
	LevelFatal
)

func (l LogLevel) string() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      *os.File
	minLevel LogLevel
	mu       sync.RWMutex
}

func makeLogFile() (*os.File, error) {
	logpath, err := os.Create("LOGS/log")
	if err != nil {
		return nil, err
	}
	return logpath, nil
}
func New(output *os.File, minLevel LogLevel) (*Logger, error) {
	output, err := makeLogFile()
	if err != nil {
		return nil, err
	}
	return &Logger{
		out:      output,
		minLevel: minLevel,
	}, nil

}

func (l *Logger) print(level LogLevel, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	temp := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.string(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}
	if level >= LevelError {
		temp.Trace = string(debug.Stack())
	}

	var report []byte

	report, err := json.Marshal(temp)
	if err != nil {
		report = []byte(LevelError.string() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.out.Write(append(report, '\n'))

}

func (l *Logger) Write(message []byte) (int, error) {
	return l.print(LevelError, string(message), nil)
}
