package service

import (
	"context"

	"github.com/MaksimovDenis/pvz_grpc/internal/models"
)

type PVZService interface {
	GetPVZ(ctx context.Context) ([]models.PVZ, error)
}
