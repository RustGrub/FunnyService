package clhouse

import (
	"context"
	"github.com/RustGrub/FunnyGoService/sql"
)

type Database struct {
	Conn Db
}

func New(c Db) *Database {
	return &Database{Conn: c}
}

func (db *Database) Exec(ctx context.Context, sql string, arguments ...any) (err error) {
	_, err = db.Conn.ExecContext(ctx, sql, arguments)
	return err
}
func (db *Database) Query(ctx context.Context, sql string, args ...any) (sql.Rows, error) {
	return db.Conn.QueryContext(ctx, sql, args)
}
func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) sql.Row {
	return db.Conn.QueryRowContext(ctx, sql, args)
}
