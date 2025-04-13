package repository

import (
	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/rs/zerolog"
)

type Repository struct {
	Authorization
	PVZ
	Receptions
	Products
}

func NewRepository(db db.Client, log zerolog.Logger) *Repository {
	return &Repository{
		Authorization: newAuthRepository(db, log),
		PVZ:           newPVZRepository(db, log),
		Receptions:    newReceptionsRepository(db, log),
		Products:      newProductsRepository(db, log),
	}
}
