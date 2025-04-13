package repository

import (
	"context"
	"time"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PVZ interface {
	CreatePVZ(ctx context.Context, newPVZ models.PVZReq) (models.PVZRes, error)
	GetFullPVZInfo(ctx context.Context, params models.GetPVZReq, limit, offset int) ([]models.FullPVZRes, error)
}

type PVZRepo struct {
	db  db.Client
	log zerolog.Logger
}

func newPVZRepository(db db.Client, log zerolog.Logger) *PVZRepo {
	return &PVZRepo{
		db:  db,
		log: log,
	}
}

func (pvz *PVZRepo) CreatePVZ(ctx context.Context, newPVZ models.PVZReq) (models.PVZRes, error) {
	var res models.PVZRes

	builder := squirrel.Insert("pvz").
		PlaceholderFormat(squirrel.Dollar).
		Columns("user_id", "city", "created_at").
		Values(newPVZ.User_id, newPVZ.City, newPVZ.RegistrationDate).
		Suffix("RETURNING id, city, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		pvz.log.Error().Err(err).Msg("CreatePVZ: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "pvz_repository.CreatePVZ",
		QueryRow: query,
	}

	err = pvz.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.City, &res.RegistrationDate)
	if err != nil {
		pvz.log.Error().Err(err).Msg("CreatePVZ: failed to execute query")
		return res, err
	}

	return res, nil

}

func (pvz *PVZRepo) GetPVZ(ctx context.Context, params models.GetPVZReq, limit, offset int) (
	[]models.PVZRes, error) {
	var res []models.PVZRes

	builder := squirrel.Select("id", "city", "created_at").
		PlaceholderFormat(squirrel.Dollar).
		From("pvz").
		Where(squirrel.Expr("created_at BETWEEN ? AND ?", params.StartDate, params.EndTime)).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := builder.ToSql()
	if err != nil {
		pvz.log.Error().Err(err).Msg("GetPVZ: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "pvz_repository.GetPVZ",
		QueryRow: query,
	}

	err = pvz.db.DB().ScanAllContext(ctx, &res, queryStruct, args...)
	if err != nil {
		pvz.log.Error().Err(err).Msg("GetPVZ: failed to scan rows")
		return nil, err
	}

	return res, nil
}

func (pvz *PVZRepo) GetFullPVZInfo(ctx context.Context, params models.GetPVZReq, limit, offset int) ([]models.FullPVZRes, error) {
	type flatRow struct {
		PVZID            uuid.UUID  `db:"pvz_id"`
		City             *string    `db:"city"`
		PVZCreatedAt     *time.Time `db:"pvz_created_at"`
		ReceptionID      uuid.UUID  `db:"reception_id"`
		ReceptionStatus  *string    `db:"reception_status"`
		ReceptionCreated *time.Time `db:"reception_created"`
		ReceptionClosed  *time.Time `db:"reception_closed"`
		ProductID        uuid.UUID  `db:"product_id"`
		ProductType      *string    `db:"product_type"`
		ProductCreatedAt *time.Time `db:"product_created"`
	}

	var rows []flatRow

	builder := squirrel.Select(
		"pvz.id AS pvz_id",
		"pvz.city",
		"pvz.created_at AS pvz_created_at",
		"r.id AS reception_id",
		"r.status AS reception_status",
		"r.created_at AS reception_created",
		"p.id AS product_id",
		"p.product_type",
		"p.created_at AS product_created",
	).
		From("pvz").
		LeftJoin("receptions r ON r.pvz_id = pvz.id").
		LeftJoin("products p ON p.reception_id = r.id").
		Where(squirrel.Expr("pvz.created_at BETWEEN ? AND ?", params.StartDate, params.EndTime)).
		OrderBy("pvz.created_at").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		pvz.log.Error().Err(err).Msg("GetFullPVZInfo: failed to build SQL")
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "pvz_repository.GetFullPVZInfo",
		QueryRow: query,
	}

	err = pvz.db.DB().ScanAllContext(ctx, &rows, queryStruct, args...)
	if err != nil {
		pvz.log.Error().Err(err).Msg("GetFullPVZInfo: failed to scan rows")
		return nil, err
	}

	pvzMap := make(map[uuid.UUID]*models.FullPVZRes)

	for _, row := range rows {
		pvzEntry, ok := pvzMap[row.PVZID]
		if !ok {
			id := row.PVZID
			pvzEntry = &models.FullPVZRes{
				Id:               &id,
				City:             derefString(row.City),
				RegistrationDate: row.PVZCreatedAt,
			}
			pvzMap[row.PVZID] = pvzEntry
		}

		var reception *models.ReceptionRes
		for i := range pvzEntry.Receptions {
			if pvzEntry.Receptions[i].Id == row.ReceptionID {
				reception = &pvzEntry.Receptions[i]
				break
			}
		}
		if reception == nil && row.ReceptionID != uuid.Nil {
			pvzEntry.Receptions = append(pvzEntry.Receptions, models.ReceptionRes{
				Id:        row.ReceptionID,
				Status:    derefString(row.ReceptionStatus),
				CreatedAt: derefTime(row.ReceptionCreated),
			})
			reception = &pvzEntry.Receptions[len(pvzEntry.Receptions)-1]
		}

		if reception != nil && row.ProductID != uuid.Nil {
			reception.Products = append(reception.Products, models.ProductRes{
				Id:          row.ProductID,
				ProductType: derefString(row.ProductType),
				CreatedAt:   derefTime(row.ProductCreatedAt),
			})
		}
	}

	var result []models.FullPVZRes
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return result, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
