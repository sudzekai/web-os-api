package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var CFG = Config{}

type Config struct {
	Logging  Logging
	WebHost  WebHost
	Database Database
	Console  Console
}

type WebHost struct {
	Host string
	Port int
}

type Logging struct {
	LogLevel string
}

type Console struct {
	IsCliEnabled bool
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
		return fmt.Errorf("Ошибка чтения файла конфигурации: Файл %s не существует", fileName)
	}

	bytes, err := os.ReadFile(fileName)

	if err != nil {
		return fmt.Errorf("Ошибка чтения файла конфигурации: %s", err.Error())
	}

	var cfg Config

	err = json.Unmarshal(bytes, &cfg)

	if err != nil {
		return fmt.Errorf("Ошибка десериализации конфигурации: %s", err.Error())
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
		return fmt.Errorf("Ошибка десериализации конфигурации: файл пуст")
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("Ошибка десериализации конфигурации: Database.Host должен быть указан")
	}
	if cfg.Database.Port == 0 {
		return fmt.Errorf("Ошибка десериализации конфигурации: Database.Port должен быть указан")
	}
	if cfg.Database.Schema == "" {
		return fmt.Errorf("Ошибка десериализации конфигурации: Database.Schema должен быть указан")
	}

	if cfg.Logging.LogLevel == "" {
		return fmt.Errorf("Ошибка десериализации конфигурации: Logging.LogLevel должен быть указан")
	}

	if cfg.WebHost.Host == "" {
		return fmt.Errorf("Ошибка десериализации конфигурации: WebHost.Host должен быть указан")
	}
	if cfg.WebHost.Port == 0 {
		return fmt.Errorf("Ошибка десериализации конфигурации: WebHost.Port должен быть указан")
	}

	return nil
}
