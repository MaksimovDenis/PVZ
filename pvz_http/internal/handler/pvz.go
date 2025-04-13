package handler

import (
	"net/http"
	"time"

	"github.com/MaksimovDenis/avito_pvz/internal/models"
	oapi "github.com/MaksimovDenis/avito_pvz/pkg/protocol"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/gin-gonic/gin"
)

func (hdl *Handler) GetPvz(ctx *gin.Context, params oapi.GetPvzParams) {
	if params.Limit == nil {
		limit := 10
		params.Limit = &limit
	}

	if params.Page == nil {
		page := 1
		params.Page = &page
	}

	if params.StartDate == nil {
		start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		params.StartDate = &start
	}

	if params.EndDate == nil {
		end := time.Now()
		params.EndDate = &end
	}

	queryParams := models.GetPVZReq{
		StartDate: *params.StartDate,
		EndTime:   *params.EndDate,
		Limit:     int(*params.Limit),
		Page:      int(*params.Page),
	}

	res, err := hdl.appService.PVZ.GetPVZ(ctx, queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (hdl *Handler) PostPvz(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	if !ok {
		hdl.log.Error().Msg("user claims not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})

		return
	}

	if !adminRole(claims.(*token.UserClaims).Role) {
		hdl.log.Error().Msg("user is not moderator")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "у пользователя нет прав"})

		return
	}

	var req oapi.PVZ

	if err := ctx.BindJSON(&req); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	newPVZ := models.PVZReq{
		User_id:          claims.(*token.UserClaims).ID,
		City:             string(req.City),
		RegistrationDate: req.RegistrationDate,
	}

	createdPVZ, err := hdl.appService.PVZ.CreatePVZ(ctx, newPVZ)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := oapi.PVZ{
		Id:               createdPVZ.Id,
		City:             oapi.PVZCity(createdPVZ.City),
		RegistrationDate: createdPVZ.RegistrationDate,
	}

	ctx.JSON(http.StatusOK, res)
}

func adminRole(role string) bool {
	if role == "moderator" {
		return true
	}

	return false
}

func employeeRole(role string) bool {
	if role == "employee" {
		return true
	}

	return false
}
