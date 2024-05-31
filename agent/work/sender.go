package work

import (
	"agent/config"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ResultData struct {
	Id     uuid.UUID `json:"id"`     // Id выражения
	Result float64   `json:"result"` // Результат выражения
}

// SendResult посылает результат оркестратору.
// `attempts` - кол-во попыток отправки результата заново при ошибке 500.
func SendResult(result ResultData, attempts uint8) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.BackendUrl+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
		return nil
	case 404:
		return errors.New("task not found")
	case 422:
		return errors.New("invalid data")
	case 500:
		if attempts == 0 {
			// Пытаемся рекурсивно отправить результат заново
			time.Sleep(time.Second)
			return SendResult(result, attempts-1)
		}
		return errors.New("orchestrator internal error")
	default:
		return errors.New("unknown error")
	}
}
