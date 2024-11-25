package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Routes) version(api *gin.RouterGroup) {
	app := api.Group("/app")
	{
		app.GET("/version", r.getVersion)
	}
}

type appVerion struct {
	Version    string `json:"version"`
	MinVersion string `json:"min_version"`
}

// @Summary Version
// @Description App version information
// @Tags version
// @Produce  json
// @Success 200 {object} appVerion
// @Router /app/version [get]
func (r *Routes) getVersion(ctx *gin.Context) {
	v := appVerion{
		Version:    "Version",
		MinVersion: "MinVersion",
	}
	ctx.JSON(http.StatusOK, v)
}
