package main

import (
	executorAbstractions "github.com/sudzekai/web-os-api/executor/abstractions"
	logging "github.com/sudzekai/web-os-api/logging/core"
	"github.com/sudzekai/web-os-api/modules/files-module/internal/controllers"
	"github.com/sudzekai/web-os-api/server/abstractions"
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
	srv abstractions.IServer,
	conf logging.LoggingConfiguration,
	executor executorAbstractions.IExecutor,
) error {
	logging.Configuration = conf

	log := logging.NewLogger("files-module:initialization")
	log.LogInformation("Запущена инициализация модуля...")

	controllers.FilesController.Connect(srv)
	return nil
}

var Module FilesModule
