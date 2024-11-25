package v1

import (
	"encoding/json"
	"net/http"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Routes) initVehicleRoutes(api *gin.RouterGroup) {
	vehicle := api.Group("/vehicle").Use(authMiddleware(r.jwtMaker))
	{
		vehicle.POST("/create", r.createVehicle)
		vehicle.POST("/delete", r.deleteVehicle)
		vehicle.PUT("/update", r.updateVehicle)
	}
}

type vehicleCreateRequest struct {
	Number  string `json:"number" binding:"required"`
	Type    string `json:"type,omitempty"` //TODO which fields are required?
	Country string `json:"country,omitempty"`
	Brand   string `json:"brand,omitempty"`
	Region  int    `json:"region,omitempty"`
}

// @Summary Create vehicle
// @Security ApiKeyAuth
// @Description Create vehicle
// @Tags vehicle
// @Accept json
// @Produce json
// @Param input body vehicleCreateRequest true "Brand, Region, Country fields are optional"
// @Success 200 {object} db.Vehicle
// @Failure 400,401,500 {object} httpServer.Error
// @Router /vehicle/create [post]
func (r *Routes) createVehicle(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("createVehicle", "getUserId(ctx)").
			Msg(err.Error())
		return
	}
	//TODO Validation for vehicle fields
	var reqBody vehicleCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("createVehicle", "vehicleCreateRequest").
			Msg(err.Error())
		return
	}
	id := uuid.NewString()
	arg := db.CreateVehicleParams{
		ID:      id,
		UserID:  userID,
		Type:    reqBody.Type,
		Brand:   reqBody.Brand,
		Number:  reqBody.Number,
		Region:  int32(reqBody.Region),
		Country: reqBody.Country,
	}
	vehicle, err := r.service.IVehicle.CreateVehicle(ctx, arg)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("createVehicle", "IVehicle.CreateVehicle").
			Str("userID", userID).
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: vehicle,
	})
}

type vehicleDeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

// @Summary Delete vehicle
// @Security ApiKeyAuth
// @Description Delete vehicle by id
// @Tags vehicle
// @Accept json
// @Produce json
// @Param input body vehicleDeleteRequest true "vehicle id"
// @Success 200 {nil} data nil
// @Failure 400,401,500 {object} httpServer.Error
// @Router /vehicle/delete [post]
func (r *Routes) deleteVehicle(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("deleteVehicle", "getUserId(ctx)").
			Msg(err.Error())
		return
	}

	var reqBody vehicleDeleteRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("deleteVehicle", "vehicleDeleteRequest").
			Msg(err.Error())
		return
	}

	arg := db.DeleteVehicleParams{
		ID:     reqBody.ID,
		UserID: userID,
	}

	err = r.service.IVehicle.DeleteVehicle(ctx, arg)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("deleteVehicle", "IVehicle.DeleteVehicle").
			Str("vehicleID", reqBody.ID).
			Str("userID", userID).
			Msg(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: nil,
	})
}

type vehicleUpdateRequest struct {
	ID      string `json:"id" binding:"required"`
	Type    string `json:"type,omitempty"`
	Brand   string `json:"brand,omitempty"`
	Number  string `json:"number,omitempty"`
	Region  int    `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}

// @Summary Update vehicle
// @Security ApiKeyAuth
// @Description Update vehicle by id
// @Tags vehicle
// @Accept json
// @Produce json
// @Param input body vehicleUpdateRequest true "Only ID is required"
// @Success 200 {object} db.Vehicle
// @Failure 400,401,404,500 {object} httpServer.Error
// @Router /vehicle/update [put]
func (r *Routes) updateVehicle(ctx *gin.Context) {

	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "getUserId(ctx)").
			Msg(err.Error())
		return
	}

	var reqBody vehicleUpdateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "vehicleUpdateRequest").
			Msg(err.Error())
		return
	}
	//TODO Validation for vehicle fields
	gotArg := db.GetVehicleByIDParams{
		ID:     reqBody.ID,
		UserID: userID,
	}
	gotVehicle, err := r.service.IVehicle.GetVehicleByID(ctx, gotArg)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "IVehicle.GetVehicleByID").
			Str("vehicleID", reqBody.ID).
			Str("userID", userID).
			Msg(err.Error())
		return
	}

	vehicleBytes, err := json.Marshal(gotVehicle)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "json.Marshal(gotVehicle)").
			Msg(err.Error())
		return
	}
	requsetBytes, err := json.Marshal(reqBody)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "json.Marshal(reqBody)").
			Msg(err.Error())
		return
	}

	patchedJSON, err := jsonpatch.MergePatch(vehicleBytes, requsetBytes)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "jsonpatch.MergePatch(vehicleBytes,requsetBytes)").
			Msg(err.Error())
		return
	}

	var updatedVehicle db.Vehicle
	err = json.Unmarshal(patchedJSON, &updatedVehicle)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "json.Unmarshal(patchedJSON, &updatedVehicle)").
			Msg(err.Error())
		return
	}

	arg := db.UpdateVehicleParams{
		Type:    pgtype.Text{String: updatedVehicle.Type, Valid: true},
		Brand:   pgtype.Text{String: updatedVehicle.Brand, Valid: true},
		Country: pgtype.Text{String: updatedVehicle.Country, Valid: true},
		Number:  pgtype.Text{String: updatedVehicle.Number, Valid: true},
		Region:  pgtype.Int4{Int32: updatedVehicle.Region, Valid: true},
		ID:      reqBody.ID,
		UserID:  userID,
	}
	vehicle, err := r.service.IVehicle.UpdateVehicle(ctx, arg)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateVehicle", "IVehicle.UpdateVehicle").
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: vehicle,
	})
}
