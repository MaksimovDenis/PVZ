package models

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	Id               uuid.UUID `json:"id" db:"id"`
	RegistrationData time.Time `json:"created_at" db:"created_at"`
	City             string    `json:"city" db:"city"`
}
