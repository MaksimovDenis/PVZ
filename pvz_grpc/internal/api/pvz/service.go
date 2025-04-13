package loader

import (
	"github.com/MaksimovDenis/pvz_grpc/internal/service"
	desc "github.com/MaksimovDenis/pvz_grpc/pkg/pvz_v1"
	"github.com/rs/zerolog"
)

type Implementation struct {
	desc.UnimplementedPVZServiceServer
	pvzSecrvice service.PVZService
	log         zerolog.Logger
}

func NewImplementation(pvzSecrvice service.PVZService, log zerolog.Logger) *Implementation {
	return &Implementation{
		pvzSecrvice: pvzSecrvice,
		log:         log,
	}
}
