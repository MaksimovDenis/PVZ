package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/MaksimovDenis/avito_pvz/internal/models"
	"github.com/MaksimovDenis/avito_pvz/internal/repository"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/MaksimovDenis/avito_pvz/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	durationAccessToken time.Duration = 24 * time.Hour
)

var invalidCharsRegex = regexp.MustCompile(`[\"'<>!#$%^&*()=+\[\]{}|\\/]`)

type Authorization interface {
	LoginUser(ctx context.Context, req models.LoginUserReq) (string, error)
	CreateUser(ctx context.Context, req models.CreateUserReq) (models.CreateUserRes, error)
	DummyLogin(ctx context.Context, role string) (string, error)
}

type AuthService struct {
	appRepository repository.Repository
	token         token.JWTMaker
	log           zerolog.Logger
}

func newAuthService(
	appRepository repository.Repository,
	token token.JWTMaker,
	log zerolog.Logger,
) *AuthService {
	return &AuthService{
		appRepository: appRepository,
		token:         token,
		log:           log,
	}
}

func (auth *AuthService) DummyLogin(ctx context.Context, role string) (string, error) {
	if err := validateRole(role); err != nil {
		return "", err
	}

	tokenInfo := models.User{}

	uuidModrator, err := uuid.Parse("11111111-1111-1111-1111-111111111111")
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to parse uuid")
		return "", err
	}

	uuidEmployee, err := uuid.Parse("22222222-2222-2222-2222-222222222222")
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to parse uuid")
		return "", err
	}

	if role == "moderator" {
		tokenInfo.Id = uuidModrator
		tokenInfo.Email = "moderatorDummyLogin@example.com"
		tokenInfo.Role = "moderator"
	} else {
		tokenInfo.Id = uuidEmployee
		tokenInfo.Email = "employeeDummyLogin@example.com"
		tokenInfo.Role = "employee"
	}

	return auth.generateToken(tokenInfo)
}

func (auth *AuthService) CreateUser(ctx context.Context, req models.CreateUserReq) (models.CreateUserRes, error) {
	var res models.CreateUserRes

	if err := validateData(req.Email, req.Password); err != nil {
		return res, err
	}

	if err := validateRole(req.Role); err != nil {
		return res, err
	}

	userId, err := uuid.NewRandom()
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to generate uuid")
		return res, errors.New("ошибка при создании нового пользователя")
	}

	req.Id = userId

	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to hash password")
		return res, errors.New("неверный логин или пароль")
	}

	req.Password = hashedPwd

	newUser, err := auth.appRepository.Authorization.CreateUser(ctx, req)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create new user in storage")
		return res, err
	}

	return newUser, nil
}

func (auth *AuthService) LoginUser(ctx context.Context, req models.LoginUserReq) (string, error) {
	if err := validateData(req.Email, req.Password); err != nil {
		return "", err
	}

	user, err := auth.appRepository.Authorization.LoginUser(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", errors.New("пользователя с данным email не существует")
		} else {
			auth.log.Error().Err(err).Msg("failed to get user from storage")
			return "", err
		}
	}

	if err = util.CheckPassword(req.Password, user.Password_hash); err != nil {
		auth.log.Error().Err(err).Msg("password mismatch")
		return "", errors.New("неверный логин или пароль")
	}

	tokenInfo := models.User{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	return auth.generateToken(tokenInfo)
}

func (auth *AuthService) generateToken(user models.User) (string, error) {
	accessToken, _, err := auth.token.CreateToken(user.Id, user.Email, user.Role, durationAccessToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create access token")
		return "", err
	}

	return accessToken, nil
}

func validateData(email, password string) error {
	switch {
	case email == "":
		return errors.New("укажите почту")
	case password == "":
		return errors.New("укажите пароль")
	case email == password:
		return errors.New("почта и пароль совпадают")
	default:
		return nil
	}
}

func validateRole(role string) error {
	if role != "employee" && role != "moderator" {
		return errors.New("Неверный формат роли пользователя")
	}

	return nil
}
