package module

import (
	"github.com/sudzekai/web-os-api/logging"
	"github.com/sudzekai/web-os-api/server"
)

type Module interface {
	Name() string
	Version() string
	Description() string
	Initialize(*server.Server, logging.LoggingConfiguration) error
}
