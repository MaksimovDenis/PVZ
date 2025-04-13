package main

import (
	"context"
	"log"

	"github.com/MaksimovDenis/avito_pvz/internal/app"
)

func main() {
	ctx := context.Background()

	avito_pvz, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	avito_pvz.Run()
}
