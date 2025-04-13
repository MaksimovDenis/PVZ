package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password_hash"`
	Role     string    `json:"role"`
}

type CreateUserReq struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

type CreateUserRes struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	Id            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password"`
	Role          string    `json:"role"`
}

type PVZReq struct {
	User_id          uuid.UUID  `json:"user_id"`
	City             string     `json:"city"`
	Id               *uuid.UUID `json:"id,omitempty"`
	RegistrationDate *time.Time `json:"registrationDate,omitempty"`
}

type PVZRes struct {
	City             string     `json:"city"`
	Id               *uuid.UUID `json:"id,omitempty"`
	RegistrationDate *time.Time `json:"registrationDate,omitempty"`
}

type GetPVZReq struct {
	StartDate time.Time `json:"startDate"`
	EndTime   time.Time `json:"endDate"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}

type CreateReceptionRes struct {
	Id       uuid.UUID `json:"id"`
	DateTime time.Time `json:"created_at"`
	PvzId    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"`
}

type LastReceptionRes struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type CreateProductReq struct {
	UserId      uuid.UUID `json:"user_id"`
	ReceptionId uuid.UUID `json:"reception_id"`
	PvzId       uuid.UUID `json:"pvz_id"`
	ProductType string    `json:"product_type"`
}

type CreateProductRes struct {
	Id          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"created_at"`
	ProductType string    `json:"product_type"`
	ReceptionId uuid.UUID `json:"reception_id"`
}

type ProductRes struct {
	Id          uuid.UUID `json:"id"`
	ProductType string    `json:"product_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReceptionRes struct {
	Id        uuid.UUID    `json:"id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	Products  []ProductRes `json:"products"`
}

type FullPVZRes struct {
	City             string         `json:"city"`
	Id               *uuid.UUID     `json:"id,omitempty"`
	RegistrationDate *time.Time     `json:"registration_date,omitempty"`
	Receptions       []ReceptionRes `json:"receptions"`
}
