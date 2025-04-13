package handler

import (
	"net/http"

	"github.com/MaksimovDenis/avito_pvz/internal/models"
	oapi "github.com/MaksimovDenis/avito_pvz/pkg/protocol"
	"github.com/MaksimovDenis/avito_pvz/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime/types"
)

func (hdl *Handler) PostProducts(ctx *gin.Context) {
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

	var req oapi.PostProductsJSONBody

	if err := ctx.BindJSON(&req); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})

		return
	}

	reqModel := models.CreateProductReq{
		UserId:      claims.(*token.UserClaims).ID,
		PvzId:       req.PvzId,
		ProductType: string(req.Type),
	}

	res, err := hdl.appService.Product.AddProduct(ctx, reqModel)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to add a product")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	resOapi := &oapi.Product{
		DateTime:    &res.DateTime,
		Id:          &res.Id,
		ReceptionId: res.ReceptionId,
		Type:        oapi.ProductType(res.ProductType),
	}

	ctx.JSON(http.StatusCreated, resOapi)
}

func (hdl *Handler) PostPvzPvzIdDeleteLastProduct(ctx *gin.Context, uuid types.UUID) {
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

	err := hdl.appService.Product.DeleteProductByPVZId(ctx, uuid)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to delete product")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.Status(http.StatusCreated)
}
