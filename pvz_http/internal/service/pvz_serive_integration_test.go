package service

import (
	"context"
	"math/rand"
	"testing"

	"github.com/MaksimovDenis/avito_pvz/internal/client/db/pg"
	"github.com/MaksimovDenis/avito_pvz/internal/client/db/transaction"
	"github.com/MaksimovDenis/avito_pvz/internal/metrics"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	pgcontainer "github.com/MaksimovDenis/avito_pvz/pkg/pg_container"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// ИНТЕГРАЦИОННЫЙ ТЕСТ
// Cоздает нового пользователя
// Cоздает новый ПВЗ
// Добавляет новую приёмку заказов
// Добавляет 50 товаров в рамках текущей приёмки заказов
// Закрывает приёмку заказов
func TestIntegration(t *testing.T) {
	ctx := context.Background()

	port := "5993"

	cli, containerID, err := pgcontainer.SetupPostgresContainer(ctx, port)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
		require.NoError(t, err)
	}()

	conStr := "postgres://admin:admin@localhost:" + port + "/testDB?sslmode=disable"

	clientDb, err := pg.New(ctx, conStr)
	if err != nil {
		t.Fatal(err)
	}
	defer clientDb.Close()

	var log zerolog.Logger

	var token token.JWTMaker

	repo := repository.NewRepository(clientDb, log)
	txManager := transaction.NewTransactionsManager(clientDb.DB())
	metrics := metrics.New()

	svc := NewService(*repo, clientDb, token, log, txManager, metrics)

	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	newUserInfo := models.CreateUserReq{
		Id:       userId,
		Email:    "admin@mail.ru",
		Password: "admin",
		Role:     "employee",
	}

	newUser, err := svc.Authorization.CreateUser(ctx, newUserInfo)
	require.NoError(t, err)

	newPVZInfo := models.PVZReq{
		User_id: newUser.Id,
		City:    "Санкт-Петербург",
	}

	newPVZ, err := svc.PVZ.CreatePVZ(ctx, newPVZInfo)
	require.NoError(t, err)

	pvzId := *newPVZ.Id

	newReception, err := svc.Reception.CreateReception(ctx, newUser.Id, pvzId)
	require.NoError(t, err)

	newProduct := models.CreateProductReq{
		UserId:      newUser.Id,
		PvzId:       pvzId,
		ReceptionId: newReception.Id,
	}

	for i := 0; i < 50; i++ {
		newProduct.ProductType = getRandomProduct()
		_, err := svc.Product.AddProduct(ctx, newProduct)
		require.NoError(t, err)
	}

	_, err = svc.Reception.CloseReceptionByPVZId(ctx, pvzId)
	require.NoError(t, err)
}

func getRandomProduct() string {
	products := []string{"одежда", "электроника", "обувь"}
	return products[rand.Intn(len(products))]
}
