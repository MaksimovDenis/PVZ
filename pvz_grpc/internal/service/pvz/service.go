package pvz

import (
	"github.com/MaksimovDenis/pvz_grpc/internal/repository"
	"github.com/MaksimovDenis/pvz_grpc/internal/service"
	"github.com/rs/zerolog"
)

type serv struct {
	pvzRepository repository.PVZRepository
	log           zerolog.Logger
}

func NewService(
	PVZRepository repository.PVZRepository,
	log zerolog.Logger,
) service.PVZService {
	return &serv{
		pvzRepository: PVZRepository,
		log:           log,
	}
}
