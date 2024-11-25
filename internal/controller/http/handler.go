package http

import (
	"net/http"

	"github.com/ShamilKhal/shgo/config"
	v1 "github.com/ShamilKhal/shgo/internal/controller/http/v1/httpRouter"
	"github.com/ShamilKhal/shgo/internal/controller/http/v1/ws"
	"github.com/ShamilKhal/shgo/internal/service"

	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service  *service.Service
	jwtMaker jwt.Maker
	config   *config.Config
}

func NewHandler(service *service.Service, jwtMaker jwt.Maker, config *config.Config) *Handler {

	return &Handler{
		jwtMaker: jwtMaker,
		service:  service,
		config:   config,
	}

}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	routesV1 := v1.NewRoutes(h.service, h.jwtMaker, h.config)
	wsChat := ws.NewWsChat(h.jwtMaker, h.service)

	{
		router.GET("/status", h.getStatus)
		routesV1.Init(router)
		wsChat.Init(router)
	}
}

func (h *Handler) getStatus(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "ok")
}
