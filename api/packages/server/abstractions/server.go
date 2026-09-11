package abstractions

import (
	"github.com/sudzekai/web-os-api/logging/abstractions"
	"github.com/sudzekai/web-os-api/server/types"
)

type IServer interface {
	Start() error
	Stop() error
	WaitForShutdown() error

	IsListening() bool

	AddMiddleware(string, types.Middleware) error

	AddHandler(string, types.HandlerFunc) error

	AddProtectedHandler(string, types.HandlerFunc, []string) error

	SetLoggingProvider(abstractions.ILogger)

	SetJWTMiddleware(types.JwtMiddleware)

	SetResultFilter(types.ResultFilter)

	GetEndpoints() []string
	GetMiddlewares() []string
	GetStats() *types.ServerStats
}
