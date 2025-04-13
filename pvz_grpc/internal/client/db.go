package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Handler func(ctx context.Context) error

type Client interface {
	DB() DB
	Close() error
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, quer Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, quer Query, args ...interface{}) error
}

type QueryExecer interface {
	ExecContext(ctx context.Context, quer Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, quer Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, quer Query, args ...interface{}) pgx.Row
}

type Query struct {
	Name     string
	QueryRow string
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Pinger
	Close()
}
