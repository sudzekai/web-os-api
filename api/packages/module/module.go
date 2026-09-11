package module

import (
	logging "github.com/sudzekai/web-os-api/logging/core"
	"github.com/sudzekai/web-os-api/server/abstractions"
)

type Module interface {
	Name() string
	Version() string
	Description() string
	Initialize(abstractions.IServer, logging.LoggingConfiguration) error
}
