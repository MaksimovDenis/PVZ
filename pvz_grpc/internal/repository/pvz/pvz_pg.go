package pvz

import (
	"context"

	db "github.com/MaksimovDenis/pvz_grpc/internal/client"
	"github.com/MaksimovDenis/pvz_grpc/internal/models"
	"github.com/MaksimovDenis/pvz_grpc/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
)

type repo struct {
	db  db.Client
	log zerolog.Logger
}

func NewRepository(db db.Client, log zerolog.Logger) repository.PVZRepository {
	return &repo{
		db:  db,
		log: log,
	}
}

func (rep *repo) GetPVZ(ctx context.Context) ([]models.PVZ, error) {
	var res []models.PVZ

	builder := squirrel.Select("id", "created_at", "city").
		PlaceholderFormat(squirrel.Dollar).
		From("pvz")

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("Get: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "pvz_repository.GetPVZ",
		QueryRow: query,
	}

	err = rep.db.DB().ScanAllContext(ctx, &res, queryStruct, args...)
	if err != nil {
		rep.log.Error().Err(err).Msg("GetPVZ: failed to scan rows")
		return nil, err
	}

	return res, nil
}
