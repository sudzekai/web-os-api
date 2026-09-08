package logging

import (
	"fmt"
	"os"
	"strings"

	pkgLogging "github.com/sudzekai/golang-logging"
)

var LoggerFactory = pkgLogging.NewLoggerFactory(os.Stdout)

func SetMinLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug", "dbug":
		LoggerFactory.SetMinLevel(pkgLogging.Debug)

	case "information", "info":
		LoggerFactory.SetMinLevel(pkgLogging.Information)

	case "warning", "warn":
		LoggerFactory.SetMinLevel(pkgLogging.Warning)

	case "error", "err":
		LoggerFactory.SetMinLevel(pkgLogging.Error)

	case "critical", "crit":
		LoggerFactory.SetMinLevel(pkgLogging.Critical)

	default:
		return fmt.Errorf("неизвестный уровень логирования: %s", level)
	}

	return nil
}
