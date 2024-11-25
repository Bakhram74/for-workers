package ws

import (
	"net/http"

	"sync"

	"github.com/ShamilKhal/shgo/internal/service"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080": // for example
		return true
	default:
		return true //TODO false
	}
}

type WSChat struct {
	jwtMaker jwt.Maker
	service  *service.Service
	clients  clientList
	sync.RWMutex
}

func NewWsChat(jwtMaker jwt.Maker, service *service.Service) *WSChat {
	return &WSChat{
		jwtMaker: jwtMaker,
		service:  service,
		clients:  make(clientList),
	}
}

func (ws *WSChat) Init(api *gin.Engine) {
	v1 := api.Group("/v1")
	{
		v1.GET("/ws", ws.ServeWS())
	}
}
