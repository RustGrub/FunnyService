package models

import (
	"encoding/json"
	"time"
)

type Good struct {
	GoodID      int       `json:"id" db:"id"`
	ProjectID   int       `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreateDt    time.Time `json:"createdAt" db:"created_at"`
}

func (g Good) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

type GoodAsLog struct {
	GoodID      int       `json:"id" db:"Id"`
	ProjectID   int       `json:"projectId" db:"ProjectId"`
	Name        string    `json:"name" db:"Name"`
	Description *string   `json:"description" db:"Description"`
	Priority    int       `json:"priority" db:"Priority"`
	Removed     bool      `json:"removed" db:"Removed"`
	CreateDt    time.Time `json:"createdAt" db:"EventTime"`
}

func (g GoodAsLog) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

type GoodsListWithMeta struct {
	Meta  Meta   `json:"meta"`
	Goods []Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type Reprioritize struct {
	GoodID    int `json:"id" db:"id"`
	ProjectID int `json:"projectId" db:"project_id"`
	Priority  int `json:"newPriority" db:"priority"`
}

type ReprioritizeResponse struct {
	Priorities []ResetPriority `json:"priorities"`
}

type ResetPriority struct {
	GoodID   int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}

type BadResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}
