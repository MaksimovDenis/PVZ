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

type Reception interface {
	CreateReception(ctx context.Context, userId, pvzId uuid.UUID) (models.CreateReceptionRes, error)
	CloseReceptionByPVZId(ctx context.Context, pvzId uuid.UUID) (models.CreateReceptionRes, error)
}

type ReceptionService struct {
	appRepository repository.Repository
	token         token.JWTMaker
	log           zerolog.Logger
	txManager     db.TxManager
	metrics       *metrics.Metrics
}

func newReceptionService(
	appRepository repository.Repository,
	token token.JWTMaker,
	log zerolog.Logger,
	txManager db.TxManager,
	metrics *metrics.Metrics,
) *ReceptionService {
	return &ReceptionService{
		appRepository: appRepository,
		token:         token,
		log:           log,
		txManager:     txManager,
		metrics:       metrics,
	}
}

func (rec *ReceptionService) CreateReception(ctx context.Context, userId, pvzId uuid.UUID) (models.CreateReceptionRes, error) {
	var res models.CreateReceptionRes

	err := rec.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		recepRes, errTx := rec.appRepository.Receptions.GetLastReceptionByPVZId(ctx, pvzId)
		if errTx != nil {
			return errors.New("неверный id ПВЗ")
		}

		if recepRes.Status == "in_progress" {
			return errors.New("невозможно начать новую приёмку товаров, пока не будет закрыта текущая")
		}

		res, errTx = rec.appRepository.Receptions.CreateReception(ctx, userId, pvzId)
		if errTx != nil {
			return errors.New("неверный запрос или есть незакрытая приемка")
		}

		return nil
	})

	if err != nil {
		return res, err
	}

	rec.metrics.ReceptionCountTotal.Inc()

	return res, nil
}

func (rec *ReceptionService) CloseReceptionByPVZId(ctx context.Context, pvzId uuid.UUID) (models.CreateReceptionRes, error) {
	var res models.CreateReceptionRes

	err := rec.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		recepRes, errTx := rec.appRepository.Receptions.GetLastReceptionByPVZId(ctx, pvzId)
		if errTx != nil {
			return errors.New("неверный id ПВЗ")
		}

		if recepRes.Status != "in_progress" {
			return errors.New("данная приёмка уже закрыта")
		}

		res, errTx = rec.appRepository.Receptions.CloseReceptionById(ctx, recepRes.Id)
		if errTx != nil {
			return errors.New("неверный запрос или приемка уже закрыта")
		}

		return nil
	})

	if err != nil {
		return res, err
	}

	return res, nil
}
