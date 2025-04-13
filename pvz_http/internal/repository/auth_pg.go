package repository

import (
	"context"
	"strings"

	db "github.com/MaksimovDenis/avito_pvz/internal/client"
	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.CreateUserReq) (models.CreateUserRes, error)
	LoginUser(ctx context.Context, req models.LoginUserReq) (models.LoginUserRes, error)
}

type AuthRepo struct {
	db  db.Client
	log zerolog.Logger
}

func newAuthRepository(db db.Client, log zerolog.Logger) *AuthRepo {
	return &AuthRepo{
		db:  db,
		log: log,
	}
}

func (arp *AuthRepo) CreateUser(ctx context.Context, user models.CreateUserReq) (models.CreateUserRes, error) {
	var res models.CreateUserRes

	builder := squirrel.Insert("users").
		PlaceholderFormat(squirrel.Dollar).
		Columns("id", "email", "password_hash", "role").
		Values(user.Id, user.Email, user.Password, user.Role).
		Suffix("RETURNING id, email, role")

	query, args, err := builder.ToSql()
	if err != nil {
		arp.log.Error().Err(err).Msg("CreateUser: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "auth_repository.CreateUser",
		QueryRow: query,
	}

	err = arp.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.Email, &res.Role)
	if err != nil {
		arp.log.Error().Err(err).Msg("CreateUser: failed to execute query")
		return res, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return res, nil
}

func (arp *AuthRepo) LoginUser(ctx context.Context, req models.LoginUserReq) (models.LoginUserRes, error) {
	var res models.LoginUserRes

	builder := squirrel.Select("id", "email", "password_hash", "role").
		PlaceholderFormat(squirrel.Dollar).
		From("users").
		Where(squirrel.Eq{"email": req.Email})

	query, args, err := builder.ToSql()
	if err != nil {
		arp.log.Error().Err(err).Msg("GetUser: failed to build SQL query")
		return res, err
	}

	queryStruct := db.Query{
		Name:     "auth_repository.LoginUser",
		QueryRow: query,
	}

	err = arp.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.Email, &res.Password_hash, &res.Role)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		arp.log.Warn().Str("email", req.Email).Msg("LoginUser: user not found")

		return res, status.Errorf(codes.NotFound, "User not found")
	} else if err != nil {
		arp.log.Error().Err(err).Msg("LoginUser: failed to execute query")

		return res, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}
