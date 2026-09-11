package logging

import (
	"fmt"
	"strings"
	"time"
)

type Logger struct {
	category string
}

func NewLogger(category string) Logger {
	return Logger{
		category: category,
	}
}

func (log *Logger) Log(logEntry LogEntry) {
	if !log.IsEnabled(logEntry.Level) {
		return
	}

	var levelStr string

	messageBytes := []byte("\r\033[K")

	if Configuration.GetCliSymbolStatus() {
		Configuration.writer.Write(messageBytes)
	}

	messageBytes = []byte(fmt.Sprintf(
		"\033[90m[%s]\033[0m %s \033[90m%s\033[0m\n",
		logEntry.TimeStamp.Format("15:04:05"),
		levelStr,
		log.category,
	))

	Configuration.writer.Write(messageBytes)

	for line := range strings.SplitSeq(logEntry.Message, "\n") {
		messageBytes = []byte(fmt.Sprintf("           %s\n", line))
		Configuration.writer.Write(messageBytes)
	}

	if Configuration.GetCliSymbolStatus() {
		messageBytes = []byte("> ")
		Configuration.writer.Write(messageBytes)
	}
}

func (log *Logger) LogDebug(format string, args ...any) {
	log.Log(makeEntry(Debug, log.category, format, args...))
}

func (log *Logger) LogInformation(format string, args ...any) {
	log.Log(makeEntry(Information, log.category, format, args...))
}

func (log *Logger) LogWarning(format string, args ...any) {
	log.Log(makeEntry(Warning, log.category, format, args...))
}

func (log *Logger) LogError(format string, args ...any) {
	log.Log(makeEntry(Error, log.category, format, args...))
}

func (log *Logger) LogCritical(format string, args ...any) {
	log.Log(makeEntry(Critical, log.category, format, args...))
} 

func makeEntry(level LogLevel, category string, format string, args ...any) LogEntry {
	return LogEntry{
		TimeStamp: time.Now(),
		Level:     level,
		Category:  category,
		Message:   fmt.Sprintf(format, args...),
	}
}

func getLevelStr(level LogLevel) string {
	switch level {
	case Debug:
		return "DBUG"

	case Information:
		return "\033[34mINFO\033[0m"

	case Warning:
		return "\033[33mWARN\033[0m"

	case Error:
		return "\033[31mERR \033[0m"

	case Critical:
		return "\033[37;41mCRIT\033[0m"
	}
}

func (log *Logger) IsEnabled(level LogLevel) bool {
	return level >= Configuration.level
}
