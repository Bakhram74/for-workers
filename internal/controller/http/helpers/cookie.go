package helpers

import (
	"net/http"

	"github.com/ShamilKhal/shgo/config"
	"github.com/gin-gonic/gin"
)

func AuthCookie(ctx *gin.Context, config *config.Config, accessToken, refreshToken string)  {

	ctx.SetSameSite(http.SameSiteNoneMode)

	ctx.SetCookie(
		config.Cookie.AccessName,
		accessToken,
		config.Cookie.AccessTtl,
		config.Cookie.AuthcookiePath,
		config.Cookie.AuthcookieDomain,
		true,
		true,
	)

	ctx.SetCookie(
		config.Cookie.RefreshName,
		refreshToken,
		config.Cookie.RefreshTtl,
		config.Cookie.AuthcookiePath,
		config.Cookie.AuthcookieDomain,
		true,
		true,
	)
}
