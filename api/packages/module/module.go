package module

import (
	executorAbstractions "github.com/sudzekai/web-os-api/executor/abstractions"
	logging "github.com/sudzekai/web-os-api/logging/core"
	serverAbstractions "github.com/sudzekai/web-os-api/server/abstractions"
)

type Module interface {
	Name() string
	Version() string
	Description() string
	Initialize(
		server serverAbstractions.IServer,
		loggingConf logging.LoggingConfiguration,
		executor executorAbstractions.IExecutor) error
}
