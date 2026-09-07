package logging

import (
	"fmt"
	"strings"
	"time"
)

// ============= CONFIGURATION =================

type loggerConfig struct {
	MinLevel LogLevel
}

var config = loggerConfig{MinLevel: Information}

func SetMinLevel(level LogLevel) {
	config.MinLevel = level
}

// ============= LOGGER ========================

type Logger struct {
	Category string
}

func NewLogger(category string) Logger {
	return Logger{
		Category: category,
	}
}

func (log *Logger) Log(level LogLevel, message string, args ...any) {
	if !isEnabled(level) {
		return
	}

	var levelStr string

	switch level {
	case Debug:
		levelStr = "DBUG"

	case Information:
		levelStr = "\033[34mINFO\033[0m"

	case Warning:
		levelStr = "\033[33mWARN\033[0m"

	case Error:
		levelStr = "\033[31mERR \033[0m"

	case Critical:
		levelStr = "\033[37;41mCRIT\033[0m"
	}

	fmt.Printf("\033[90m[%s]\033[0m %s \033[90m%s\033[0m\n", time.Now().Format("15:04:05"), levelStr, log.Category)
	message = fmt.Sprintf(message, args...)

	for line := range strings.SplitSeq(message, "\n") {
		fmt.Printf("           %s\n", line) // проблема со встроенным pad если передается fmt.sprintf или другая форматированная строка
	}
}

func (log *Logger) LogDebug(message string, args ...any) {
	log.Log(Debug, message, args...)
}

func (log *Logger) LogInformation(message string, args ...any) {
	log.Log(Information, message, args...)
}

func (log *Logger) LogWarning(message string, args ...any) {
	log.Log(Warning, message, args...)
}

func (log *Logger) LogError(message string, args ...any) {
	log.Log(Error, message, args...)
}

func (log *Logger) LogCritical(message string, args ...any) {
	log.Log(Critical, message, args...)
}

// ============= UTILITIES =====================

func isEnabled(level LogLevel) bool {
	var result bool = false

	switch config.MinLevel {
	case Debug:
		if level == Debug || level == Information || level == Warning || level == Error || level == Critical {
			return true
		}
	case Information:
		if level == Information || level == Warning || level == Error || level == Critical {
			return true
		}
	case Warning:
		if level == Warning || level == Error || level == Critical {
			return true
		}
	case Error:
		if level == Error || level == Critical {
			return true
		}
	case Critical:
		if level == Critical {
			return true
		}
	}

	return result
}
