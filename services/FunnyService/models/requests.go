package models

// Для каждого запроса можно использовать good из models, но я предпочту под каждый запрос свое

// CreateGoodRequest В хендлере используется для ожидания name, из query приходит projId и включается для записи в бд
type CreateGoodRequest struct {
	Name      string `json:"name" db:"name"`
	ProjectID int    `json:"projectId" db:"project_id"`
}

type UpdateGoodRequest struct {
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	ProjectID   int     `json:"projectId" db:"project_id"`
	GoodID      int     `json:"id" db:"id"`
}

type RemoveGoodRequest struct {
	ProjectID int  `json:"projectId" db:"project_id"`
	GoodID    int  `json:"id" db:"id"`
	Removed   bool `json:"removed" db:"removed"`
}
