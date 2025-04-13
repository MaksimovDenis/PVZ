package handler

import (
	"net/http"

	oapi "github.com/MaksimovDenis/avito_pvz/pkg/protocol"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime/types"
)

func (hdl *Handler) PostReceptions(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	if !ok {
		hdl.log.Error().Msg("user claims not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})

		return
	}

	if !employeeRole(claims.(*token.UserClaims).Role) {
		hdl.log.Error().Msg("user is not employee")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "у пользователя нет прав"})

		return
	}

	var req oapi.PostReceptionsJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	userId := claims.(*token.UserClaims).ID
	pvzId := req.PvzId

	res, err := hdl.appService.Reception.CreateReception(ctx, userId, pvzId)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to create new reception")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	resOapi := &oapi.Reception{
		DateTime: res.DateTime,
		Id:       &res.Id,
		PvzId:    res.PvzId,
		Status:   oapi.ReceptionStatus(res.Status),
	}

	ctx.JSON(http.StatusOK, resOapi)
}

func (hdl *Handler) PostPvzPvzIdCloseLastReception(ctx *gin.Context, uuid types.UUID) {
	claims, ok := ctx.Get("user")
	if !ok {
		hdl.log.Error().Msg("user claims not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})

		return
	}

	if !employeeRole(claims.(*token.UserClaims).Role) {
		hdl.log.Error().Msg("user is not employee")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "у пользователя нет прав"})

		return
	}

	res, err := hdl.appService.Reception.CloseReceptionByPVZId(ctx, uuid)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to close reception")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	resOapi := &oapi.Reception{
		DateTime: res.DateTime,
		Id:       &res.Id,
		PvzId:    res.PvzId,
		Status:   oapi.ReceptionStatus(res.Status),
	}

	ctx.JSON(http.StatusOK, resOapi)
}
