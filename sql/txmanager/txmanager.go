package txmanager

import (
	"context"
	"errors"
	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbTagTX = iota

type Manager struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Manager {
	return &Manager{db: db}
}

func (m *Manager) Exec(ctx context.Context, sql string, arguments ...any) (err error) {
	_, err = m.getDB(ctx).Exec(ctx, sql, arguments...)
	return err
}

func (m *Manager) Query(ctx context.Context, sql string, args ...any) (sql.Rows, error) {
	return m.getDB(ctx).Query(ctx, sql, args...)
}

func (m *Manager) QueryRow(ctx context.Context, sql string, args ...any) sql.Row {
	return m.getDB(ctx).QueryRow(ctx, sql, args...)
}

func (m *Manager) Begin(ctx context.Context, transactional TransactionFn) (err error) {
	tx, err := m.db.Begin(ctx)
	if tx == nil {
		return errors.New(consts.ErrNoDatabase)
	}
	defer func() {
		if err != nil {
			rErr := tx.Rollback(ctx)
			//subsidary.NewLogUsingContext(ctx, "Rolling back tx")
			if rErr != nil && errors.Is(rErr, pgx.ErrTxClosed) {
				// если транзакция уже закрыта или отменена не надо писать о проблеме роллбека
			}
		}
	}()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, dbTagTX, tx)

	err = transactional(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *Manager) getDB(ctx context.Context) DB {
	tx, ok := ctx.Value(dbTagTX).(pgx.Tx)
	if ok {
		return tx
	}
	return m.db
}
