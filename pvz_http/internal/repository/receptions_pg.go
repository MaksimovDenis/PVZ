package repository

import (
	"context"
	"strings"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Receptions interface {
	CreateReception(ctx context.Context, userId, pvzId uuid.UUID) (models.CreateReceptionRes, error)
	GetLastReceptionByPVZId(ctx context.Context, pvzId uuid.UUID) (models.LastReceptionRes, error)
	CloseReceptionById(ctx context.Context, receptionId uuid.UUID) (models.CreateReceptionRes, error)
}

type ReceptionsRepo struct {
	db  db.Client
	log zerolog.Logger
}

func newReceptionsRepository(db db.Client, log zerolog.Logger) *ReceptionsRepo {
	return &ReceptionsRepo{
		db:  db,
		log: log,
	}
}

func (rec *ReceptionsRepo) CreateReception(ctx context.Context, userId, pvzId uuid.UUID) (models.CreateReceptionRes, error) {
	var res models.CreateReceptionRes

	builder := squirrel.Insert("receptions").
		PlaceholderFormat(squirrel.Dollar).
		Columns("user_id", "pvz_id", "status").
		Values(userId.String(), pvzId.String(), "in_progress").
		Suffix("RETURNING id, created_at, pvz_id, status")

	query, args, err := builder.ToSql()
	if err != nil {
		rec.log.Error().Err(err).Msg("CreateReception: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "receptions_repository.CreateReception",
		QueryRow: query,
	}

	err = rec.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.DateTime, &res.PvzId, &res.Status)
	if err != nil {
		rec.log.Error().Err(err).Msg("CreateReception: failed to execute query")
		return res, err
	}

	return res, nil
}

func (rec *ReceptionsRepo) GetLastReceptionByPVZId(ctx context.Context, pvzId uuid.UUID) (models.LastReceptionRes, error) {
	var res models.LastReceptionRes

	builder := squirrel.Select("id", "status").
		PlaceholderFormat(squirrel.Dollar).
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzId}).
		OrderBy("created_at DESC").
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		rec.log.Error().Err(err).Msg("GetLastReceptionByPVZId: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "receptions_repository.GetLastReceptionByPVZId",
		QueryRow: query,
	}

	err = rec.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.Status)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return res, nil
	} else if err != nil {
		rec.log.Error().Err(err).Msg("GetLastReceptionByPVZId: failed to execute query")
		return res, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}

func (rec *ReceptionsRepo) CloseReceptionById(ctx context.Context, receptionId uuid.UUID) (models.CreateReceptionRes, error) {
	var res models.CreateReceptionRes

	builder := squirrel.Update("receptions").
		PlaceholderFormat(squirrel.Dollar).
		Set("status", "close").
		Where(squirrel.Eq{"id": receptionId}).
		Suffix("RETURNING id, created_at, pvz_id, status")

	query, args, err := builder.ToSql()
	if err != nil {
		rec.log.Error().Err(err).Msg("CloseReceptionById: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "receptions_repository.CloseReceptionById",
		QueryRow: query,
	}

	err = rec.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.DateTime, &res.PvzId, &res.Status)
	if err != nil {
		rec.log.Error().Err(err).Msg("CloseReceptionById: failed to execute query")
		return res, err
	}

	return res, nil
}
