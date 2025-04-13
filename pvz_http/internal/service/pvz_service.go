package service

import (
	"context"
	"errors"
	"time"

	"github.com/MaksimovDenis/avito_pvz/internal/metrics"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/rs/zerolog"
)

type PVZ interface {
	CreatePVZ(ctx context.Context, req models.PVZReq) (models.PVZRes, error)
	GetPVZ(ctx context.Context, req models.GetPVZReq) ([]models.FullPVZRes, error)
}

type PVZService struct {
	appRepository repository.Repository
	token         token.JWTMaker
	log           zerolog.Logger
	metrics       *metrics.Metrics
}

func newPVZService(
	appRepository repository.Repository,
	token token.JWTMaker,
	log zerolog.Logger,
	metrics *metrics.Metrics,
) *PVZService {
	return &PVZService{
		appRepository: appRepository,
		token:         token,
		log:           log,
		metrics:       metrics,
	}
}

func (pvz *PVZService) CreatePVZ(ctx context.Context, newPVZ models.PVZReq) (models.PVZRes, error) {
	var res models.PVZRes

	if err := validateCity(newPVZ.City); err != nil {
		return res, err
	}

	if newPVZ.RegistrationDate == nil {
		now := time.Now()
		newPVZ.RegistrationDate = &now
	}

	res, err := pvz.appRepository.PVZ.CreatePVZ(ctx, newPVZ)
	if err != nil {
		return res, errors.New("ошибка при создании нового ПВЗ")
	}

	pvz.metrics.PvzCountTotal.Inc()

	return res, nil
}

func (pvz *PVZService) GetPVZ(ctx context.Context, req models.GetPVZReq) ([]models.FullPVZRes, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	res, err := pvz.appRepository.PVZ.GetFullPVZInfo(ctx, req, req.Limit, offset)
	if err != nil {
		return res, errors.New("ошибка при получении списка ПВЗ")
	}

	return res, nil
}

func validateCity(city string) error {
	if city != "Москва" && city != "Казань" &&
		city != "Санкт-Петербург" {
		return errors.New("в данном городе пока нет доступных ПВЗ")
	}

	return nil
}
