package app

import (
	"context"
	"net"

	"github.com/MaksimovDenis/pvz_grpc/internal/closer"
	"github.com/MaksimovDenis/pvz_grpc/internal/config"
	desc "github.com/MaksimovDenis/pvz_grpc/pkg/pvz_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/rs/zerolog/log"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCServer()
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(app.grpcServer)

	desc.RegisterPVZServiceServer(app.grpcServer, app.serviceProvider.LoadImpl(ctx))

	return nil
}

func (app *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", app.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", app.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = app.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
