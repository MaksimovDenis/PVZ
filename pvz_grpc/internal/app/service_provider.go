package app

import (
	"context"
	"os"

	pvz "github.com/MaksimovDenis/pvz_grpc/internal/api/pvz"
	db "github.com/MaksimovDenis/pvz_grpc/internal/client"
	"github.com/MaksimovDenis/pvz_grpc/internal/client/db/pg"
	"github.com/MaksimovDenis/pvz_grpc/internal/closer"
	"github.com/MaksimovDenis/pvz_grpc/internal/config"
	"github.com/MaksimovDenis/pvz_grpc/internal/repository"
	pvzRepository "github.com/MaksimovDenis/pvz_grpc/internal/repository/pvz"
	pvzService "github.com/MaksimovDenis/pvz_grpc/internal/service/pvz"

	"github.com/MaksimovDenis/pvz_grpc/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient      db.Client
	pvzRepository repository.PVZRepository

	pvzService service.PVZService

	log zerolog.Logger

	pvzImpl *pvz.Implementation
}

func newServiceProvider() *serviceProvider {
	srv := &serviceProvider{}
	srv.log = srv.initLogger()

	return srv
}

func (srv *serviceProvider) initLogger() zerolog.Logger {
	logFile, err := os.OpenFile("./internal/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}

	logLevel, err := zerolog.ParseLevel("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	multiWriter := zerolog.MultiLevelWriter(os.Stdout, logFile)

	logger := zerolog.New(multiWriter).Level(logLevel).With().Timestamp().Logger()

	return logger
}

func (srv *serviceProvider) PGConfig() config.PGConfig {
	if srv.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get pg config")
		}

		srv.pgConfig = cfg
	}

	return srv.pgConfig
}

func (srv *serviceProvider) GRPCConfig() config.GRPCConfig {
	if srv.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get server config")
		}

		srv.grpcConfig = cfg
	}

	return srv.grpcConfig
}

func (srv *serviceProvider) DBClient(ctx context.Context) db.Client {
	if srv.dbClient == nil {
		client, err := pg.New(ctx, srv.PGConfig().DSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create db client")
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("ping error")
		}

		closer.Add(func() error {
			client.Close()
			return nil
		})

		srv.dbClient = client
	}

	return srv.dbClient
}

func (srv *serviceProvider) LoaderRepository(ctx context.Context) repository.PVZRepository {
	if srv.pvzRepository == nil {
		srv.pvzRepository = pvzRepository.NewRepository(srv.DBClient(ctx), srv.log)
	}

	return srv.pvzRepository
}

func (srv *serviceProvider) LoaderService(ctx context.Context) service.PVZService {
	if srv.pvzService == nil {
		srv.pvzService = pvzService.NewService(
			srv.LoaderRepository(ctx),
			srv.log)
	}

	return srv.pvzService
}

func (srv *serviceProvider) LoadImpl(ctx context.Context) *pvz.Implementation {
	if srv.pvzImpl == nil {
		srv.pvzImpl = pvz.NewImplementation(srv.LoaderService(ctx), srv.log)
	}

	return srv.pvzImpl
}
