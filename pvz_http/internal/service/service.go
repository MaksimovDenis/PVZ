package service

import (
	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/metrics"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/rs/zerolog"
)

type Service struct {
	Authorization
	PVZ
	Reception
	Product
}

func NewService(repos repository.Repository,
	client db.Client,
	token token.JWTMaker,
	log zerolog.Logger,
	txManager db.TxManager,
	metrics *metrics.Metrics) *Service {
	return &Service{
		Authorization: newAuthService(repos, token, log),
		PVZ:           newPVZService(repos, token, log, metrics),
		Reception:     newReceptionService(repos, token, log, txManager, metrics),
		Product:       newProductService(repos, token, log, txManager, metrics),
	}
}
