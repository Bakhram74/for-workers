package v1

import (
	"github.com/ShamilKhal/shgo/config"
	"github.com/ShamilKhal/shgo/internal/service"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	service  *service.Service
	jwtMaker jwt.Maker
	config   *config.Config
}

func NewRoutes(service *service.Service, jwtMaker jwt.Maker, config *config.Config) *Routes {
	return &Routes{
		service:  service,
		jwtMaker: jwtMaker,
		config:   config,
	}
}

func (r *Routes) Init(api *gin.Engine) {
	v1 := api.Group("/v1")
	{
		r.initAuthRoutes(v1)
		r.initUserRoutes(v1)
		r.initVehicleRoutes(v1)
		r.initVehiclePageRoutes(v1)
		r.initChatRoutes(v1)
		r.version(v1)
	}
}
