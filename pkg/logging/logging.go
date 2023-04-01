package logging

import (
	"encoding/json"
	"fmt"
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
	Out      *os.File
	minLevel LogLevel
	mu       sync.RWMutex
}

func makeLogFile() (*os.File, error) {
	logFile, err := os.OpenFile("LOG.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("something about the file", err)
		return nil, err
	}
	// defer logFile.Close()
	return logFile, nil
}

func NewLogger() (*Logger, error) {
	output, err := makeLogFile()
	if err != nil {
		return nil, err
	}
	return &Logger{
		Out:      output,
		minLevel: LevelInfo,
	}, nil

}

func (l *Logger) LogInfo(message, source string) {
	l.log(LevelInfo, message, source)
}

func (l *Logger) LogError(err error, source string) {
	l.log(LevelError, err.Error(), source)
}

func (l *Logger) LogFatal(err error, source string) {
	l.log(LevelFatal, err.Error(), source)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, message, source string) (int, error) {
	temp := struct {
		Level   string `json:"level"`
		Source  string `json:"source"`
		Time    string `json:"time"`
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	}{
		Level:   level.String(),
		Time:    time.Now().UTC().Format(time.RFC3339),
		Source:  source,
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
	return l.Out.Write(append(report, '\n'))

}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.log(LevelError, string(message), "LOGEER")
}
