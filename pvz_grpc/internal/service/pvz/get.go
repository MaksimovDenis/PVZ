package pvz

import (
	"context"
	"errors"

	"github.com/MaksimovDenis/pvz_grpc/internal/models"
)

func (srv *serv) GetPVZ(ctx context.Context) ([]models.PVZ, error) {
	pvzList, err := srv.pvzRepository.GetPVZ(ctx)
	if err != nil {
		return nil, errors.New("ошибка при получении списка ПВЗ")
	}

	return pvzList, nil
}
