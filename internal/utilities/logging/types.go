package logging

type LogLevel string

const (
	Debug       LogLevel = "Debug"
	Information LogLevel = "Information"
	Warning     LogLevel = "Warning"
	Error       LogLevel = "Error"
	Critical    LogLevel = "Critical"
)
