package sql

import (
	"context"
)

type Database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (err error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}
