package calculate

import "github.com/google/uuid"

type DataRequest struct {
	Expression string `json:"expression"` // Строка с выражением
}

type DataResponse struct {
	Id uuid.UUID `json:"id"` // Идентификатор выражения (UUID)
}
