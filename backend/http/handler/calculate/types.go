package calculate

type DataRequest struct {
	Expression string `json:"expression"` // Строка с выражением
}

type DataResponse struct {
	Id int64 `json:"id"` // Идентификатор выражения (UUID)
}
