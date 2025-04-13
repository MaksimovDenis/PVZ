package repository

import (
	"context"

	"github.com/MaksimovDenis/pvz_grpc/internal/models"
)

type PVZRepository interface {
	GetPVZ(ctx context.Context) ([]models.PVZ, error)
}
