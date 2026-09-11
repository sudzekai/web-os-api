package main

import (
	"github.com/sudzekai/web-os-api/logging"
	"github.com/sudzekai/web-os-api/server"
)

type FilesModule struct{}

func (FilesModule) Name() string {
	return "Files Module"
}

func (FilesModule) Version() string {
	return "0.0.1"
}

func (FilesModule) Description() string {
	return ""
}

func (FilesModule) Initialize(
	srv *server.Server,
	conf logging.LoggingConfiguration,
) error {
	logging.Configuration = conf

	log := logging.NewLogger("files-module:initialization")
	log.LogInformation("Запущена инициализация модуля...")

	return nil
}

var Module FilesModule
