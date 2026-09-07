package main

import (
	"os"

	"github.com/sudzekai/web-os-api/internal/args"
	"github.com/sudzekai/web-os-api/internal/config"
	"github.com/sudzekai/web-os-api/internal/server"
	"github.com/sudzekai/web-os-api/internal/utilities/logging"
)

func main() {
	logging.SetMinLevel(logging.Debug)

	log := logging.NewLogger("main")

	args.LoadArguments(os.Args)

	err := config.LoadConfig("appsettings.json")

	if err != nil {
		log.LogCritical("%s", err.Error())
		return
	}

	log.LogInformation(
		"Система запускается с следующей конфигурацией:\nЛогирование:\n\tУровень: %s\nХостинг:\n\tХост: %s\n\tПорт: %d\nБаза данных:\n\tХост: %s\n\tПорт: %d\n\tСхема: %s",
		config.CFG.Logging.LogLevel,
		config.CFG.WebHost.Host,
		config.CFG.WebHost.Port,
		config.CFG.Database.Host,
		config.CFG.Database.Port,
		config.CFG.Database.Schema,
	)

	server := server.NewServer(config.CFG.WebHost.Host, config.CFG.WebHost.Port)

	err = server.Start()
}
