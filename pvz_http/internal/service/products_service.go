package service

import (
	"context"
	"errors"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/metrics"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Product interface {
	AddProduct(ctx context.Context, req models.CreateProductReq) (models.CreateProductRes, error)
	DeleteProductByPVZId(ctx context.Context, pvzId uuid.UUID) error
}

type ProductService struct {
	appRepository repository.Repository
	token         token.JWTMaker
	log           zerolog.Logger
	txManager     db.TxManager
	metrics       *metrics.Metrics
}

func newProductService(
	appRepository repository.Repository,
	token token.JWTMaker,
	log zerolog.Logger,
	txManager db.TxManager,
	metrics *metrics.Metrics,
) *ProductService {
	return &ProductService{
		appRepository: appRepository,
		token:         token,
		log:           log,
		txManager:     txManager,
		metrics:       metrics,
	}
}

func (prd *ProductService) AddProduct(ctx context.Context, req models.CreateProductReq) (models.CreateProductRes, error) {
	var res models.CreateProductRes

	if err := validateProductType(req.ProductType); err != nil {
		return res, err
	}

	err := prd.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		recepRes, errTx := prd.appRepository.Receptions.GetLastReceptionByPVZId(ctx, req.PvzId)
		if errTx != nil {
			return errors.New("неверный запрос или нет активной приемки")
		}

		if recepRes.Status != "in_progress" {
			return errors.New("неверный запрос или нет активной приемки")
		}

		req.ReceptionId = recepRes.Id

		res, errTx = prd.appRepository.Products.AddProduct(ctx, req)
		if errTx != nil {
			return errors.New("неверный запрос или нет активной приемки")
		}

		return nil
	})

	if err != nil {
		return res, err
	}

	prd.metrics.ProductsCountTotal.Inc()

	return res, nil
}

func (prd *ProductService) DeleteProductByPVZId(ctx context.Context, pvzId uuid.UUID) error {
	err := prd.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		recepRes, errTx := prd.appRepository.Receptions.GetLastReceptionByPVZId(ctx, pvzId)
		if errTx != nil {
			return errors.New("небходимо создать новую приёмку товаров для данного ПВЗ")
		}

		if recepRes.Status == "close" {
			return errors.New("приёмка товаров в данном ПВЗ закрыта, необходимо открыть новую")
		}

		productId, errTx := prd.appRepository.Products.GetLastProductIdByReceptionId(ctx, recepRes.Id)
		if errTx != nil {
			return errors.New("в рамках текущей приёмки нет товаров для удаления")
		}

		errTx = prd.appRepository.Products.DeleteProduct(ctx, productId)
		if errTx != nil {
			return errors.New("неверный запрос, нет активной приемки или нет товаров для удаления")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func validateProductType(product string) error {
	if product != "электроника" && product != "одежда" &&
		product != "обувь" {
		return errors.New("данный тип товара не поддерживается")
	}

	return nil
}
