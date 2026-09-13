package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	executor "github.com/sudzekai/web-os-api/executor/abstractions"
	logging "github.com/sudzekai/web-os-api/logging/abstractions"
	"github.com/sudzekai/web-os-api/module"
	"github.com/sudzekai/web-os-api/server/abstractions"
	server "github.com/sudzekai/web-os-api/server/abstractions"
)

type ModulesLoader struct {
	modules        map[string]string
	loggingFactory logging.ILoggerFactory
	logger         logging.ILogger
	executor       executor.IExecutor
	server         server.IServer
}

func NewModulesLoader(loggerFactory logging.ILoggerFactory) *ModulesLoader {
	return &ModulesLoader{
		modules: make(map[string]string),
		logger:  loggerFactory.NewLogger("modules-loader"),
	}
}

var modules map[string]string

func (loader *ModulesLoader) Load(path string, srv abstractions.IServer) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("загрузка модуля %q: %w", path, err)
	}

	symbol, err := p.Lookup("Module")
	if err != nil {
		return fmt.Errorf(
			"модуль %q не экспортирует Module: %w",
			path,
			err,
		)
	}

	mod, ok := symbol.(module.Module)
	if !ok {
		return fmt.Errorf(
			"модуль %q имеет неправильный тип Module",
			path,
		)
	}

	if err := mod.Initialize(srv, loggingConfiguration); err != nil {
		return fmt.Errorf(
			"инициализация модуля %q: %w",
			mod.Name(),
			err,
		)
	}

	modules[mod.Name()] = mod.Version()

	log.LogDebug("Модуль %s загружен", mod.Name())

	return nil
}

func LoadModules(srv abstractions.IServer, loggingConfiguration logging.LoggingConfiguration) error {
	modules = make(map[string]string)

	log := logging.NewLogger("modules")
	log.LogInformation("Начата загрузка модулей...")

	files, err := os.ReadDir("./modules")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			log.LogInformation("Загрузка модуля: %s...", file.Name())
			err = Load(filepath.Join("./modules", file.Name()), srv, loggingConfiguration)
			if err != nil {
				log.LogError("Ошибка загрузки модуля: %s. Модуль пропущен...", err.Error())
			}
		}
	}

	LogLoadedModules()

	return nil
}

func LogLoadedModules() {
	log := logging.NewLogger("modules:dictionary")

	var message []string = make([]string, 0)
	for name, ver := range modules {
		message = append(message, fmt.Sprintf("%s (v%s)", name, ver))
	}

	log.LogInformation("Загруженные модули:\n%s", strings.Join(message, "\n"))

}
