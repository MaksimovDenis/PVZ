package main

import (
	"context"
	"log"

	"github.com/MaksimovDenis/pvz_grpc/internal/app"
)

func main() {
	ctx := context.Background()

	pvz, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = pvz.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
