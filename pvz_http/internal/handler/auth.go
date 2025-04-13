package handler

import (
	"net/http"

	"github.com/MaksimovDenis/avito_pvz/internal/models"
	oapi "github.com/MaksimovDenis/avito_pvz/pkg/protocol"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime/types"
)

func (hdl *Handler) PostDummyLogin(ctx *gin.Context) {
	var dummyLogin oapi.PostDummyLoginJSONBody

	if err := ctx.BindJSON(&dummyLogin); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	token, err := hdl.appService.DummyLogin(ctx, string(dummyLogin.Role))
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to auth user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (hdl *Handler) PostRegister(ctx *gin.Context) {
	var retigterReq oapi.PostRegisterJSONBody

	if err := ctx.BindJSON(&retigterReq); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	modelReq := models.CreateUserReq{
		Email:    string(retigterReq.Email),
		Password: retigterReq.Password,
		Role:     string(retigterReq.Role),
	}

	user, err := hdl.appService.Authorization.CreateUser(ctx, modelReq)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to create user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	res := &oapi.User{
		Id:    &user.Id,
		Email: types.Email(user.Email),
		Role:  oapi.UserRole(user.Role),
	}

	ctx.JSON(http.StatusOK, res)
}

func (hdl *Handler) PostLogin(ctx *gin.Context) {
	var loginReq oapi.PostLoginJSONBody

	if err := ctx.BindJSON(&loginReq); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	modelsReq := models.LoginUserReq{
		Email:    string(loginReq.Email),
		Password: string(loginReq.Password),
	}

	token, err := hdl.appService.LoginUser(ctx, modelsReq)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to auth user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, token)
}
