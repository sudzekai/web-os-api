package config

import (
	"encoding/json"
	"os"

	"github.com/sudzekai/web-os-api/internal/errors"
	"github.com/sudzekai/web-os-api/internal/utilities/logging"
)

var CFG = Config{}

type Config struct {
	Logging  Logging
	WebHost  WebHost
	Database Database
}

type WebHost struct {
	Host string
	Port int
}

type Logging struct {
	LogLevel logging.LogLevel
}

type Database struct {
	Host     string
	Port     int
	User     string
	Schema   string
	Password string
}

func LoadConfig(fileName string) error {
	_, err := os.Stat(fileName)

	if err != nil {
		return errors.Config.ReadFile("Файл %s не существует", fileName)
	}

	bytes, err := os.ReadFile(fileName)

	if err != nil {
		return errors.Config.ReadFile("%s", err.Error())
	}

	var cfg Config

	err = json.Unmarshal(bytes, &cfg)

	if err != nil {
		return errors.Config.UnmarshalFile("%s", err.Error())
	}

	err = validate(cfg)

	if err != nil {
		return err
	}

	CFG = cfg

	return nil
}

func validate(cfg Config) error {
	if cfg == CFG {
		return errors.Config.UnmarshalFile("Конфигурация пуста")
	}

	if cfg.Database.Host == "" {
		return errors.Config.UnmarshalFile("Database.Host должен быть указан")
	}
	if cfg.Database.Port == 0 {
		return errors.Config.UnmarshalFile("Database.Port должен быть указан")
	}
	if cfg.Database.Schema == "" {
		return errors.Config.UnmarshalFile("Database.Schema должен быть указан")
	}

	if cfg.Logging.LogLevel == "" {
		return errors.Config.UnmarshalFile("Logging.LogLevel должен быть указан")
	}

	if cfg.WebHost.Host == "" {
		return errors.Config.UnmarshalFile("WebHost.Host должен быть указан")
	}
	if cfg.WebHost.Port == 0 {
		return errors.Config.UnmarshalFile("WebHost.Port должен быть указан")
	}

	return nil
}
