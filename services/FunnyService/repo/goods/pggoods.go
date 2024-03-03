package goods

import (
	"context"
	"errors"
	"fmt"
	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/RustGrub/FunnyGoService/sql/txmanager"
	"github.com/jackc/pgx/v5"
)

type Db struct {
	db     txmanager.DB
	logger logger.Logger
}

func NewPgRepo(db txmanager.DB, logger logger.Logger) *Db {
	return &Db{db: db, logger: logger}
}

// Тут, скорее всего, стоит сделать возврат указателей на структуры, но пока так

const createQuery = "INSERT INTO goods (project_id, name) VALUES ($1,$2) RETURNING id, project_id, name, description, priority, removed, created_at;"

func (r *Db) Create(ctx context.Context, req models.CreateGoodRequest) (res *models.Good, err error) {
	var good models.Good
	err = r.db.QueryRow(ctx, createQuery, req.ProjectID, req.Name).Scan(
		&good.GoodID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreateDt,
	)
	if err != nil {
		err = fmt.Errorf("error creating good: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	r.logger.Warning(fmt.Sprintf("ReqID: %v:Inserted good with name %v and projectID %v into db", ctx.Value(consts.ReqID).(string), good.Name, good.ProjectID))

	return &good, nil
}

const updateNameAndDescQuery = "UPDATE goods SET name=($1), description=($2) WHERE id=($3) AND project_id=($4) RETURNING id, project_id, name, description, priority, removed, created_at;"

func (r *Db) UpdateNameAndDescription(ctx context.Context, req models.UpdateGoodRequest) (res *models.Good, err error) {
	var good models.Good
	err = r.db.QueryRow(ctx, updateNameAndDescQuery, req.Name, req.Description, req.GoodID, req.ProjectID).Scan(
		&good.GoodID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreateDt,
	)
	if err != nil {
		// Если не пришли строки, то не нашли по айдишникам
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		err = fmt.Errorf("error updating good: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	r.logger.Warning(fmt.Sprintf("ReqID: %v:updated good with name %v and projectID %v into db", ctx.Value(consts.ReqID).(string), good.Name, good.ProjectID))

	return &good, nil
}

const removeQuery = "UPDATE goods SET removed=true WHERE id=($1) AND project_id=($2) AND removed=false RETURNING id, project_id, name, description, priority, removed, created_at;"

// Remove Если правильно понял, то просто помечаем removed как true, но продолжаем хранить
func (r *Db) Remove(ctx context.Context, req models.RemoveGoodRequest) (res *models.Good, err error) {
	var good models.Good
	err = r.db.QueryRow(ctx, removeQuery, req.GoodID, req.ProjectID).Scan(
		&good.GoodID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreateDt,
	)
	if err != nil {
		// Если не пришли строки, то не нашли по айдишникам, либо уже удалена
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		err = fmt.Errorf("error removing good: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	r.logger.Warning(fmt.Sprintf("ReqID: %v:removed good with id %v and projectID %v into db", ctx.Value(consts.ReqID).(string), good.Name, good.ProjectID))

	return &good, nil
}

const getListByLimitAndOffsetQuery = "SELECT * FROM goods WHERE priority >= ($1) ORDER BY priority LIMIT ($2)"

func (r *Db) GetListByLimitAndOffset(ctx context.Context, limit, offset int) (res []models.Good, removed int, err error) {
	rows, err := r.db.Query(ctx, getListByLimitAndOffsetQuery, offset, limit)
	if err != nil {
		err = fmt.Errorf("error getting list of repo: %v", err)
		r.logger.Error(err)
		return []models.Good{}, 0, err
	}

	res, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Good, error) {
		var good models.Good
		err := row.Scan(
			&good.GoodID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreateDt,
		)
		if good.Removed {
			removed++
		}
		return good, err
	})
	if err != nil {
		// Возвращаем пустой слайс без ошибки (отобразим пустым)
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Good{}, 0, nil
		}
		err = fmt.Errorf("error getting list of repo: %v", err)
		r.logger.Error(err)
		return []models.Good{}, 0, err
	}
	return
}

const reprioritizeGoodsUpQuery = "UPDATE goods SET priority = priority - 1 WHERE priority > ($1) AND priority <=($2) RETURNING id, project_id, name, description, priority, removed, created_at;"
const reprioritizeGoodsDownQuery = "UPDATE goods SET priority = priority + 1 WHERE priority >= ($1) AND priority < ($2) RETURNING id, project_id, name, description, priority, removed, created_at;"

func (r *Db) Reprioritize(ctx context.Context, reqP, goodP int, up bool) (res []models.Good, err error) {
	var rows pgx.Rows

	if up {
		rows, err = r.db.Query(ctx, reprioritizeGoodsUpQuery, goodP, reqP)
	} else {
		rows, err = r.db.Query(ctx, reprioritizeGoodsDownQuery, reqP, goodP)
	}
	if err != nil {
		err = fmt.Errorf("error getting list of repo: %v", err)
		r.logger.Error(err)
		return []models.Good{}, err
	}

	res, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Good, error) {
		var good models.Good
		err := row.Scan(
			&good.GoodID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreateDt,
		)
		return good, err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Good{}, nil
		}
		err = fmt.Errorf("error getting list of repo: %v", err)
		r.logger.Error(err)
		return []models.Good{}, err
	}
	return
}

const getGoodQuery = "SELECT * FROM goods WHERE id=($1) AND project_id=($2)"

func (r *Db) Get(ctx context.Context, goodID, projectID int) (res *models.Good, err error) {
	var good models.Good
	err = r.db.QueryRow(ctx, getGoodQuery, goodID, projectID).Scan(
		&good.GoodID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreateDt,
	)
	if err != nil {
		// Если не пришли строки, то не нашли по айдишникам
		if errors.Is(err, pgx.ErrNoRows) {
			err = errors.New(consts.ErrGoodNotFound)
		}
		err = fmt.Errorf("error getting good: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	r.logger.Warning(fmt.Sprintf("ReqID: %v:got good with name %v and projectID %v from db", ctx.Value(consts.ReqID).(string), good.Name, good.ProjectID))

	return &good, nil
}

const updatePriorityQuery = "UPDATE goods SET priority=($1) WHERE id=($2) AND project_id=($3) RETURNING id, project_id, name, description, priority, removed, created_at;"

func (r *Db) UpdatePriority(ctx context.Context, goodID, projectID, priority int) (res *models.Good, err error) {
	var good models.Good
	err = r.db.QueryRow(ctx, updatePriorityQuery, priority, goodID, projectID).Scan(
		&good.GoodID,
		&good.ProjectID,
		&good.Name,
		&good.Description,
		&good.Priority,
		&good.Removed,
		&good.CreateDt,
	)
	if err != nil {
		err = fmt.Errorf("error updating good: %v", err)
		r.logger.Error(err)
		return nil, err
	}
	r.logger.Warning(fmt.Sprintf("ReqID: %v:updated good with name %v and projectID %v into db", ctx.Value(consts.ReqID).(string), good.Name, good.ProjectID))

	return &good, nil
}
