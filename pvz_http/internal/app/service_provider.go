package app

import (
	"context"
	"os"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/client/db/pg"
	"github.com/MaksimovDenis/avito_pvz/internal/client/db/transaction"
	"github.com/MaksimovDenis/avito_pvz/internal/closer"
	"github.com/MaksimovDenis/avito_pvz/internal/config"
	"github.com/MaksimovDenis/avito_pvz/internal/handler"
	"github.com/MaksimovDenis/avito_pvz/internal/metrics"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	"github.com/MaksimovDenis/avito_pvz/internal/service"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	serverConfig config.ServerConfig
	tokenConfig  config.TokenConfig

	dbClient      db.Client
	txManager     db.TxManager
	appRepository *repository.Repository

	appService *service.Service

	handler *handler.Handler

	tokenMaker *token.JWTMaker

	log     zerolog.Logger
	metrics *metrics.Metrics
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

func (srv *serviceProvider) initMetric() *metrics.Metrics {
	srv.metrics = metrics.New()
	return srv.metrics
}

func (srv *serviceProvider) ServerConfig() config.ServerConfig {
	if srv.serverConfig == nil {
		cfg, err := config.NewServerConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get server config")
		}

		srv.serverConfig = cfg
	}

	return srv.serverConfig
}

func (srv *serviceProvider) TokenConfig() config.TokenConfig {
	if srv.tokenConfig == nil {
		cfg, err := config.NewSecretKey()
		if err != nil {
			log.Fatal().Err(err).Msg("failed dto get secret key config")
		}

		srv.tokenConfig = cfg
	}

	return srv.tokenConfig
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

func (srv *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if srv.txManager == nil {
		srv.txManager = transaction.NewTransactionsManager(srv.DBClient(ctx).DB())
	}

	return srv.txManager
}

func (srv *serviceProvider) TokenMaker(ctx context.Context) *token.JWTMaker {
	if srv.tokenMaker == nil {
		srv.tokenMaker = token.NewJWTMaker(
			srv.TokenConfig().SecretKey(),
		)
	}

	return srv.tokenMaker
}

func (srv *serviceProvider) AppRepository(ctx context.Context) *repository.Repository {
	if srv.appRepository == nil {
		srv.appRepository = repository.NewRepository(
			srv.DBClient(ctx),
			srv.log.With().Str("module", "repository").Logger(),
		)
	}

	return srv.appRepository
}

func (srv *serviceProvider) AppService(ctx context.Context) *service.Service {
	if srv.appService == nil {
		srv.appService = service.NewService(
			*srv.AppRepository(ctx),
			srv.DBClient(ctx),
			*srv.TokenMaker(ctx),
			srv.log.With().Str("module", "service").Logger(),
			srv.TxManager(ctx),
			srv.initMetric(),
		)
	}

	return srv.appService
}

func (srv *serviceProvider) AppHandler(ctx context.Context) *handler.Handler {
	if srv.handler == nil {
		srv.handler = handler.NewHandler(
			*srv.AppService(ctx),
			*srv.TokenMaker(ctx),
			srv.log.With().Str("module", "api").Logger(),
			srv.metrics,
		)
	}

	return srv.handler
}
