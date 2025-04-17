package logging

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

type Logger struct {
	Level   LogLevel
	LogFile *os.File
}

func openLogFile() (*os.File, error) {
	file, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (l *Logger) LogMessage(level LogLevel, message string) {
	if level <= l.Level {
		t := time.Now().Format(time.RFC3339)
		switch level {
		case LogLevelInfo:
			l.LogFile.WriteString(fmt.Sprintf("[%s] [INFO]: %s\n", t, message))
		case LogLevelWarning:
			l.LogFile.WriteString(fmt.Sprintf("[%s] [WARNING]: %s\n", t, message))
		case LogLevelError:
			l.LogFile.WriteString(fmt.Sprintf("[%s] [ERROR]: %s\n", t, message))
		case LogLevelDebug:
			l.LogFile.WriteString(fmt.Sprintf("[%s] [DEBUG]: %s\n", t, message))
		default:
			l.LogFile.WriteString(fmt.Sprintf("[%s] [UNKNOWN]: %s\n", t, message))
		}
	}
}

func NewLogger(level LogLevel) (*Logger, error) {
	file, err := openLogFile()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		Level:   level,
		LogFile: file,
	}

	return logger, nil
}
