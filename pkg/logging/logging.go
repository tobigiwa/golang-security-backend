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

func (l LogLevel) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "INVALID ERROR LEVEL"
	}
}

type Logger struct {
	out      *os.File
	minLevel LogLevel
	mu       sync.RWMutex
}

func makeLogFile() (*os.File, error) {
	logFile, err := os.OpenFile("LOGS/log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer logFile.Close()
	return logFile, nil
}

func NewLogger(minLevel LogLevel) (*Logger, error) {
	output, err := makeLogFile()
	if err != nil {
		return nil, err
	}
	return &Logger{
		out:      output,
		minLevel: minLevel,
	}, nil

}

func (l *Logger) LogInfo(message, source string) {
	l.print(LevelInfo, message, source)
}

func (l *Logger) LogError(err error, source string) {
	l.print(LevelError, err.Error(), source)
}

func (l *Logger) LogFatal(err error, source string) {
	l.print(LevelFatal, err.Error(), source)
	os.Exit(1)
}

func (l *Logger) print(level LogLevel, message, source string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	temp := struct {
		Level   string `json:"level"`
		Source  string `json:"sorce"`
		Time    string `json:"time"`
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    time.Now().UTC().Format(time.RFC3339),
		Message: message,
	}
	if level >= LevelError {
		temp.Trace = string(debug.Stack())
	}

	var report []byte

	report, err := json.Marshal(temp)
	if err != nil {
		report = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.out.Write(append(report, '\n'))

}

func (l *Logger) Write(message []byte) (int, error) {
	return l.print(LevelError, string(message), "LOGEER")
}
