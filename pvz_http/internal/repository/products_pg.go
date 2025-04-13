package repository

import (
	"context"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Products interface {
	AddProduct(ctx context.Context, req models.CreateProductReq) (models.CreateProductRes, error)
	GetLastProductIdByReceptionId(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error)
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type ProductsRepo struct {
	db  db.Client
	log zerolog.Logger
}

func newProductsRepository(db db.Client, log zerolog.Logger) *ProductsRepo {
	return &ProductsRepo{
		db:  db,
		log: log,
	}
}

func (prd *ProductsRepo) AddProduct(ctx context.Context, req models.CreateProductReq) (models.CreateProductRes, error) {
	var res models.CreateProductRes

	builder := squirrel.Insert("products").
		PlaceholderFormat(squirrel.Dollar).
		Columns("user_id", "pvz_id", "reception_id", "product_type").
		Values(req.UserId, req.PvzId, req.ReceptionId, req.ProductType).
		Suffix("RETURNING id, created_at, product_type, reception_id")

	query, args, err := builder.ToSql()
	if err != nil {
		prd.log.Error().Err(err).Msg("AddProduct: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "products_repository.AddProduct",
		QueryRow: query,
	}

	err = prd.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.DateTime, &res.ProductType, &res.ReceptionId)
	if err != nil {
		prd.log.Error().Err(err).Msg("AddProduct: failed to execute query")
		return res, err
	}

	return res, nil
}

func (prd *ProductsRepo) GetLastProductIdByReceptionId(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error) {
	var productId uuid.UUID

	builder := squirrel.Select("id").
		PlaceholderFormat(squirrel.Dollar).
		From("products").
		Where(squirrel.Eq{"reception_id": receptionId}).
		OrderBy("created_at DESC").
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		prd.log.Error().Err(err).Msg("GetLastProductIdByReceptionId: failed to build SQL query")
		return productId, err
	}

	queryStruct := db.Query{
		Name:     "products_repository.GetLastProductIdByReceptionId",
		QueryRow: query,
	}

	err = prd.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&productId)
	if err != nil {
		prd.log.Error().Err(err).Msg("GetLastProductIdByReceptionId: failed to execute query")
		return productId, err
	}

	return productId, nil
}

func (prd *ProductsRepo) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	builder := squirrel.Delete("products").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": productId})

	query, args, err := builder.ToSql()
	if err != nil {
		prd.log.Error().Err(err).Msg("DeleteProduct: failed to build SQL query")
		return err
	}

	queryStruct := db.Query{
		Name:     "products_repository.DeleteProduct",
		QueryRow: query,
	}

	_, err = prd.db.DB().ExecContext(ctx, queryStruct, args...)
	if err != nil {
		prd.log.Error().Err(err).Msg("DeleteProduct: failed to execute query")
		return err
	}

	return nil
}
