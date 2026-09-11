package server

type LoggingProvider interface {
	LogDebug(string, ...any)
	LogInformation(string, ...any)
	LogWarning(string, ...any)
	LogError(string, ...any)
	LogCritical(string, ...any)
}
