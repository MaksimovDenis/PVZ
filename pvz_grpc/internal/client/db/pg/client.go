package pg

import (
	"context"
	"time"

	db "github.com/MaksimovDenis/pvz_grpc/internal/client"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

const (
	maxRetries = 3
	retryDelay = 3 * time.Second
)

func New(ctx context.Context, dsn string) (db.Client, error) {
	var dbc *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		dbc, err = pgxpool.Connect(ctx, dsn)
		if err == nil {
			break
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, errors.Errorf("failed to connect db after %d attempts: %v", maxRetries, err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (clt *pgClient) DB() db.DB {
	return clt.masterDBC
}

func (clt *pgClient) Close() error {
	if clt.masterDBC != nil {
		clt.masterDBC.Close()
	}

	return nil
}
