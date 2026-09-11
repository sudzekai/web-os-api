package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"sudzekai/logging"
	"github.com/sudzekai/web-os-api/internal/args"
	"github.com/sudzekai/web-os-api/internal/config"
	"github.com/sudzekai/web-os-api/internal/controllers"
	"github.com/sudzekai/web-os-api/internal/server"
)

func main() {
	configureEnvironment()

	srv := createServer()

	go srv.Start()

	if config.CFG.Console.IsCliEnabled {
		readInput(srv)
		return
	}

	srv.WaitForShutdown()
}

func createServer() *server.Server {
	srv := server.NewServer(
		config.CFG.WebHost.Host,
		config.CFG.WebHost.Port,
	)

	srv.AddHandler("GET /test", controllers.TestHandler)
	srv.AddHandler("GET /test/{id}", controllers.TestHandler)

	return srv
}

func configureEnvironment() {
	configureArgs()
	applyConfiguration()
}

func configureArgs() {
	log := logging.LoggerFactory.NewLogger("main:configuration:args")

	args.LoadArguments(os.Args)

	if err := loadConfiguration(); err != nil {
		log.LogCritical("%s", err.Error())
		os.Exit(-1)
	}

	applyArgumentOverrides()
}

func loadConfiguration() error {
	if args.Args["--dev"] != "" {
		return config.LoadConfig("appsettings.development.json")
	}

	cfgFileName := args.Args["--cfg"]

	if cfgFileName != "" && cfgFileName != "val" {
		return config.LoadConfig(cfgFileName)
	}

	return config.LoadConfig("appsettings.json")
}

func applyArgumentOverrides() {
	if levelStr := args.Args["--log-level"]; levelStr != "" {
		config.CFG.Logging.LogLevel = levelStr
	}

	if args.Args["--cli"] != "" {
		config.CFG.Console.IsCliEnabled = true
	}
}

func applyConfiguration() {
	cfg := config.CFG

	if err := logging.SetMinLevel(cfg.Logging.LogLevel); err != nil {
		panic(err)
	}

	if cfg.Console.IsCliEnabled {
		logging.LoggerFactory.EnableCliSymbol()
	}

	log := logging.LoggerFactory.NewLogger("main:configuration:apply")

	logArguments(log)
	logConfiguration(log, &cfg)
}

func logArguments(log pkgLogging.Logger) {
	if len(args.Args) == 0 {
		return
	}

	arguments := make([]string, 0, len(args.Args))

	for flag, value := range args.Args {
		arguments = append(
			arguments,
			fmt.Sprintf("%s %s", flag, value),
		)
	}

	log.LogDebug(
		"Аргументы:\n%s",
		strings.Join(arguments, "\n"),
	)
}

func logConfiguration(log pkgLogging.Logger, cfg *config.Config) {
	cfg.Database.Password = "*****"
	cfg.Database.User = "*****"

	data, _ := json.MarshalIndent(cfg, "", "\t")

	log.LogDebug(
		"Конфигурация:\n%s",
		string(data),
	)
}

func readInput(srv *server.Server) {
	reader := bufio.NewReader(os.Stdin)
	log := logging.LoggerFactory.NewLogger("main:cli")

	commands := createCommands(srv, log)

	for {
		input, err := reader.ReadString('\n')

		if err != nil {
			return
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			return
		}

		if command, ok := commands[input]; ok {
			command()
			continue
		}

		log.LogInformation("Неизвестная команда")
	}
}

func createCommands(
	srv *server.Server,
	log pkgLogging.Logger,
) map[string]func() {
	return map[string]func(){
		"stop": srv.Stop,

		"start": func() {
			go srv.Start()
		},

		"restart": func() {
			srv.Stop()
			go srv.Start()
		},

		"endpoints": func() {
			logEndpoints(srv, log)
		},

		"stat": func() {
			logStats(srv, log)
		},
	}
}

func logEndpoints(srv *server.Server, log pkgLogging.Logger) {
	log.LogInformation(
		"Эндпоинты:\n\t%s",
		strings.Join(srv.GetEndpoints(), "\n\t"),
	)
}

func logStats(srv *server.Server, log pkgLogging.Logger) {
	stat := srv.GetStats()

	log.LogInformation(
		"Статистика:\n%-17s %v\n%-17s %d\n%-17s %d\n%-17s %d\n%-17s %s\n%-17s %s",
		"Запущен:", srv.IsListening,
		"Запросы:", stat.Requests.Load(),
		"Ответы:", stat.Responses.Load(),
		"Ошибки:", stat.Errors.Load(),
		"Время запуска:", stat.StartTime.Format("15:04:05"),
		"Время остановки:", stat.StopTime.Format("15:04:05"),
	)
}
