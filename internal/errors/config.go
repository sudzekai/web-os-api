package errors

import (
	"fmt"
)

type config string

const Config config = "cfg"

func (config) ReadFile(err string, args ...any) error {
	return MkError("Ошибка чтения файла конфигурации: %s", fmt.Sprintf(err, args...))
}

func (config) UnmarshalFile(err string, args ...any) error {
	return MkError("Ошибка десериализации файла конфигурации: %s", fmt.Sprintf(err, args...))
}
