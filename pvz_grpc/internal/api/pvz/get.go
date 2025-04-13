package loader

import (
	"context"

	"github.com/MaksimovDenis/pvz_grpc/internal/models"
	"github.com/MaksimovDenis/pvz_grpc/pkg/pvz_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (hdl *Implementation) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	list, err := hdl.pvzSecrvice.GetPVZ(ctx)
	if err != nil {
		hdl.log.Error().Err(err).Msgf("failed to get pvz list")
		return nil, err
	}

	return &pvz_v1.GetPVZListResponse{
		Pvzs: converterToListPVZRes(list),
	}, nil
}

func converterModelToPVZRes(data models.PVZ) *pvz_v1.PVZ {
	return &pvz_v1.PVZ{
		Id:               data.Id.String(),
		RegistrationDate: timestamppb.New(data.RegistrationData),
		City:             data.City,
	}
}

func converterToListPVZRes(list []models.PVZ) []*pvz_v1.PVZ {
	pvzList := make([]*pvz_v1.PVZ, len(list))

	for idx, value := range list {
		pvzList[idx] = converterModelToPVZRes(value)
	}

	return pvzList
}
