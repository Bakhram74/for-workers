package v1

import (
	"errors"
	"net/http"

	"github.com/ShamilKhal/shgo/internal/controller/http/flashcall"
	"github.com/ShamilKhal/shgo/internal/controller/http/helpers"
	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/ShamilKhal/shgo/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (r *Routes) initAuthRoutes(api *gin.RouterGroup) {

	auth := api.Group("/auth")
	{
		auth.POST("/login", r.login)
		auth.POST("/verify", r.verifyUser)
		auth.POST("/refresh", r.refreshToken)
	}
	auth.Use(authMiddleware(r.jwtMaker))
	auth.POST("/create", r.createUser)

}

type loginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// @Summary Login
// @Description Create user if user doesnt exist and return id.This endpoint also generate pincode and send it to user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body loginRequest true "phone number"
// @Success 200 {object} service.stickerData "sticker"
// @Failure      400,404,500  {object}  httpServer.Error
// @Router /auth/login [post]
func (r *Routes) login(ctx *gin.Context) {

	var reqBody loginRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("login", "loginRequest").
			Msg(err.Error())
		return
	}

	phone, err := utils.ValidatePhone(reqBody.Phone)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("login", "ValidatePhone").
			Str("phone", reqBody.Phone).
			Msg(err.Error())
		return
	}

	stickerData, err := r.service.IAuthorization.Login(ctx, phone)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("login", "IAuthorization.Login").
			Msg(err.Error())
		return
	}

	attempts, err := flashcall.Flashcall(ctx, r.config, r.service, phone, stickerData.Sticker)
	logger.Log.Info().Str("login", "Flashcall").Msgf("Retry attempts: %d", attempts)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("login", "Flashcall").
			Int("attempts left", attempts).
			Msg(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: stickerData,
	})
}

type verifyUserRequest struct {
	Sticker string `json:"sticker" binding:"required"`
	Pincode string `json:"pincode" binding:"required"`
}

// @Summary Verify user
// @Description Compare pincode by sticker(id) and return user. This endpoint also sets authentication cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body verifyUserRequest true "sticker & pincode"
// @Success 200 {object} db.User "Returns user object and sets authentication cookies"
// @Failure 400,404,500 {object} httpServer.Error
// @Router /auth/verify [post]
func (r *Routes) verifyUser(ctx *gin.Context) {
	var req verifyUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("verifyUser", "verifyUserRequest").
			Msg(err.Error())
		return
	}

	data, err := r.service.IAuthorization.AuthVerifyUser(ctx, req.Sticker, req.Pincode)
	if err != nil {
		if err.Error() == "redis: wrong pincode" {
			httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			logger.Log.Error().
				Str("verifyUser", "IAuthorization.AuthVerifyUser").
				Str("pincode", req.Pincode).
				Msg(err.Error())
			return
		}
		if err.Error() == "redis: key not found" {
			httpServer.ErrorResponse(ctx, http.StatusNotFound, err.Error())
			logger.Log.Error().
				Str("verifyUser", "IAuthorization.AuthVerifyUser").
				Str("key", req.Sticker).
				Msg(err.Error())
			return
		}
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("verifyUser", "IAuthorization.AuthVerifyUser").
			Msg(err.Error())
		return
	}

	helpers.AuthCookie(ctx, r.config, data.AccessToken, data.RefreshToken)

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: data.User,
	})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Refresh token
// @Description Accept refresh_token and set AccessToken and RefreshToken in cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body refreshTokenRequest true "refresh_token"
// @Success 200 {nil} data "Returns nil and sets authentication cookies"
// @Failure      400,404,500  {object}  httpServer.Error
// @Router /auth/refresh [post]
func (r *Routes) refreshToken(ctx *gin.Context) {

	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("refreshToken", "refreshTokenRequest").
			Msg(err.Error())
		return
	}

	tokens, err := r.service.IAuthorization.AuthRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) || errors.Is(err, jwt.ErrExpiredToken) {
			httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			logger.Log.Error().
				Str("refreshToken", "IAuthorization.AuthRefreshToken").
				Str("refresh_token", req.RefreshToken).
				Msg(err.Error())
			return
		}
		if err.Error() == "no rows in result set" {
			httpServer.ErrorResponse(ctx, http.StatusNotFound, err.Error())
			logger.Log.Error().
				Str("refreshToken", "IAuthorization.AuthRefreshToken").
				Str("refresh_token", req.RefreshToken).
				Msg(err.Error())
			return
		}
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("refreshToken", "IAuthorization.AuthRefreshToken").
			Msg(err.Error())
		return
	}

	helpers.AuthCookie(ctx, r.config, tokens.AccessToken, tokens.RefreshToken)

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: nil,
	})
}

type createUserrequest struct {
	Name       string `json:"name" binding:"required"`
	ImageUrl   string `json:"image_url,omitempty"`
	StatusText string `json:"status_text,omitempty"`
}

// @Summary Create user
// @Security ApiKeyAuth
// @Description Check JWT token and update user fieldsand and also sets authentication cookies
// @Tags auth
// @Accept json
// @Produce json
// @Param input body createUserrequest true " Name is required, while ImageUrl and StatusText are optional"
// @Success 200 {object} db.User
// @Failure 400,404,500 {object} httpServer.Error
// @Router /auth/create [post]
func (r *Routes) createUser(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("createUser", "getUserId(ctx)").
			Msg(err.Error())
		return
	}

	var reqBody createUserrequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		logger.Log.Error().
			Str("createUser", "createUserrequest").
			Msg(err.Error())
		return
	}

	data, err := r.service.IAuthorization.AuthCreateUser(ctx, userID, reqBody.Name, reqBody.ImageUrl, reqBody.StatusText)
	if err != nil {
		if err.Error() == "user already exists" {
			httpServer.ErrorResponse(ctx, http.StatusForbidden, err.Error())
			logger.Log.Error().
				Str("createUser", "IAuthorization.AuthCreateUser").
				Str("userID", userID).
				Msg(err.Error())
			return
		}
		if err.Error() == "no rows in result set" {
			httpServer.ErrorResponse(ctx, http.StatusNotFound, err.Error())
			logger.Log.Error().
				Str("createUser", "IAuthorization.AuthCreateUser").
				Str("userID", userID).
				Msg(err.Error())
			return
		}
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("createUser", "IAuthorization.AuthCreateUser").
			Msg(err.Error())
		return
	}

	helpers.AuthCookie(ctx, r.config, data.AccessToken, data.RefreshToken)

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: data.User,
	})
}
