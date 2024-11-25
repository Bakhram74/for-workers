package v1

import (
	"net/http"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"

	"github.com/ShamilKhal/shgo/internal/controller/http/flashcall"
	"github.com/ShamilKhal/shgo/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (r *Routes) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user").Use(authMiddleware(r.jwtMaker))
	{
		user.PUT("/update-phone", r.updateUserPhone)
		user.PUT("/verify-phone", r.verifyUserPhone)
		user.PUT("/update-img", r.updateUserImg)
		user.PUT("/update-data", r.updateUserData)
	}

}

type updateUserPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary Update phone
// @Security ApiKeyAuth
// @Description Generate pincode and send it to user
// @Tags user
// @Accept json
// @Produce json
// @Param input body updateUserPhoneRequest true "New phone number"
// @Success 200 {nil} data nil
// @Failure 400,401,500 {object} httpServer.Error
// @Router /user/update-phone [put]
func (r *Routes) updateUserPhone(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("updateUserPhone", "getUserId(ctx)").
			Msg(err.Error())
		return
	}

	var reqBody updateUserPhoneRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("updateUserPhone", "updateUserPhoneRequest").
			Msg(err.Error())
		return
	}
	phone, err := utils.ValidatePhone(reqBody.Phone)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("updateUserPhone", "ValidatePhone").
			Str("phone", reqBody.Phone).
			Msg(err.Error())
		return
	}
	attempts, err := flashcall.Flashcall(ctx, r.config, r.service, phone, userID)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateUserPhone", "Flashcall").
			Int("attempts left", attempts).
			Msg(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: nil,
	})
}

type verifyUserPhoneRequest struct {
	Pincode string `json:"pincode" binding:"required"`
}

// @Summary Verify phone
// @Security ApiKeyAuth
// @Description Compare pincode and reset phone
// @Tags user
// @Accept json
// @Produce json
// @Param input body verifyUserPhoneRequest true "pincode"
// @Success 200 {object} db.User "Returns user object"
// @Failure 400,401,404,500 {object} httpServer.Error
// @Router /user/verify-phone [put]
func (r *Routes) verifyUserPhone(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("verifyUserPhone", "getUserId(ctx)").
			Msg(err.Error())
		return
	}
	var reqBody verifyUserPhoneRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("verifyUserPhone", "verifyUserPhoneRequest").
			Msg(err.Error())
		return
	}
	user, err := r.service.IUser.UpdateUserPhone(ctx, userID, reqBody.Pincode)
	if err != nil {
		if err.Error() == "redis: wrong pincode" {
			httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			logger.Log.Error().
				Str("verifyUserPhone", "IUser.UpdateUserPhone").
				Str("pincode", reqBody.Pincode).
				Msg(err.Error())
			return
		}
		if err.Error() == "redis: key not found" {
			httpServer.ErrorResponse(ctx, http.StatusNotFound, err.Error())
			logger.Log.Error().
				Str("verifyUserPhone", "IUser.UpdateUserPhone").
				Str("key", userID).
				Msg(err.Error())
			return
		}
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("verifyUserPhone", "IUser.UpdateUserPhone").
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: user,
	})
}

type updateUserImgRequest struct {
	ImageUrl string `json:"image_url" binding:"required"`
}

// @Summary Update image
// @Security ApiKeyAuth
// @Description Update users image
// @Tags user
// @Accept json
// @Produce json
// @Param input body updateUserImgRequest true "image_url"
// @Success 200 {object} db.User "Returns user object"
// @Failure 400,401,500 {object} httpServer.Error
// @Router /user/update-img [put]
func (r *Routes) updateUserImg(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("updateUserImg", "getUserId(ctx)").
			Msg(err.Error())
		return
	}
	var reqBody updateUserImgRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("updateUserImg", "updateUserImgRequest").
			Msg(err.Error())
		return
	}
	user, err := r.service.IUser.UpdateUserImg(ctx, userID, reqBody.ImageUrl)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateUserImg", "IUser.UpdateUserImg").
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: user,
	})
}

type updateUserDataRequest struct {
	Name       string `json:"name,omitempty"`
	StatusText string `json:"status_text,omitempty"`
}

// @Summary Update  data
// @Security ApiKeyAuth
// @Description Update users name & status_text
// @Tags user
// @Accept json
// @Produce json
// @Param input body updateUserDataRequest true "users name & status_text fields are optional"
// @Success 200 {object} db.User "Returns user object"
// @Failure 400,401,500 {object} httpServer.Error
// @Router /user/update-data [put]
func (r *Routes) updateUserData(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("updateUserData", "getUserId(ctx)").
			Msg(err.Error())
		return
	}
	var reqBody updateUserDataRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("updateUserData", "updateUserDataRequest").
			Msg(err.Error())
		return
	}
	if reqBody.Name == "" && reqBody.StatusText == "" {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, "empty fields are provided")
		logger.Log.Error().
			Str("updateUserData", "name & status_text").
			Msg("empty fields are provided")
		return
	}
	user, err := r.service.IUser.UpdateUserData(ctx, userID, reqBody.Name, reqBody.StatusText)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("updateUserData", "IUser.UpdateUserData").
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: user,
	})
}
