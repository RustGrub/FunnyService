package goods

import (
	"context"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/RustGrub/FunnyGoService/sql"
)

type Db struct {
	db     sql.Database
	logger logger.Logger
}

func New(db sql.Database, logger logger.Logger) *Db {
	return &Db{db: db, logger: logger}
}

const setGoodQuery = "INSERT INTO goods (GoodId, ProjectId, Name, Description, Priority, Removed, EventTime) VALUES ($1,$2,$3,$4,$5,$6,$7)"

func (r *Db) Set(ctx context.Context, good models.GoodAsLog) error {
	err := r.db.Exec(ctx, setGoodQuery, good.GoodID, good.ProjectID, good.Name, good.Description, good.Priority, good.Removed, good.CreateDt)
	if err != nil {
		return err
	}
	return nil
}
