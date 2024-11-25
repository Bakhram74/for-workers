package v1

import (
	"net/http"

	"strconv"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (r *Routes) initVehiclePageRoutes(api *gin.RouterGroup) {
	api.POST("/vehicle/page", r.vehicleFind)
}

type vehicleRequest struct {
	Number string `json:"number" binding:"required"`
	Region int    `json:"region,omitempty"`
}

type vehicleResponse struct {
	Vehicles   []db.Vehicle `json:"vehicles"`
	Pagination `json:"pagination"`
}

type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
}

// @Summary Find vehicle
// @Description Find vehicle by number and region
// @Tags vehicle
// @Accept json
// @Produce json
// @Param			limit	query		string		true "count of vehicles"
// @Param			page	query		string		true "page of vehicles"
// @Param input body vehicleRequest true "Region is optional"
// @Success 200 {object} vehicleResponse
// @Failure 400,500 {object} httpServer.Error
// @Router /vehicle/page [post]
func (r *Routes) vehicleFind(ctx *gin.Context) {

	var reqBody vehicleRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("vehicleFind", "vehicleRequest").
			Msg(err.Error())
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("vehicleFind", "ctx.Query(page)").
			Msg(err.Error())
		return
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("vehicleFind", "ctx.Query(limit)").
			Msg(err.Error())
		return
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	offset := limit * (page - 1)

	vehicles, count, err := r.service.IVehicle.FindVehicle(ctx, reqBody.Number, reqBody.Region, limit, offset)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("vehicleFind", "IVehicle.FindVehicle").
			Str("number", reqBody.Number).
			Int("region", reqBody.Region).
			Int("limit", limit).
			Int("page", page).
			Msg(err.Error())
		return
	}
	totalVehicles := (count / limit)
	remainder := (count % limit)
	if remainder != 0 {
		totalVehicles = totalVehicles + 1
	}

	response := vehicleResponse{
		Vehicles: vehicles,
		Pagination: Pagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalVehicles,
		},
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: response,
	})
}
