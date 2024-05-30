package calculate

import (
	"backend/internal/utils"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Парсинг тела запроса
	var data DataRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Проверка, что выражение не пустое
	if data.Expression == "" {
		http.Error(w, "expression must not be empty", http.StatusUnprocessableEntity)
		return
	}

	// Проверка, что выражение не содержит запрещённые символы
	var allowedChars = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '(', ')', '/', '*', '^', '.', ' '}
	if contains, char := utils.CheckCharsInString(data.Expression, allowedChars); !contains {
		http.Error(w, "forbidden symbol found: "+string(char), http.StatusUnprocessableEntity)
		return
	}

	// TODO: Добавление в очередь
	logger.Info(data.Expression)
	id := uuid.Must(uuid.NewRandom()) // допустим генерируем id

	// Отправка ответа
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(DataResponse{Id: id})
	if _, err = w.Write(response); err != nil {
		logger.Error(err.Error())
	}
}
