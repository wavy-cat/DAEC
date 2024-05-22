package helper

import (
	"encoding/json"
	"io"
)

// ParseJSON функция для демаршализации данных JSON из io.Reader в структуру T.
func ParseJSON[T any](reader io.Reader) (T, error) {
	var data T

	body, err := io.ReadAll(reader) // Читаем данные
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data) // Демаршализируем JSON
	if err != nil {
		return data, err
	}
	return data, nil
}
