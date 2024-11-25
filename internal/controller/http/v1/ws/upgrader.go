package ws

import (
	"net/http"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (ws *WSChat) ServeWS() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.Query("token")
		payload, err := ws.jwtMaker.VerifyToken(token)
		if err != nil {
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			logger.Log.Error().
				Str("ServeWS()", "jwtMaker.VerifyToken").
				Str("token", token).
				Msg(err.Error())
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			httpServer.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			logger.Log.Error().
				Str("ServeWS()", "upgrader.Upgrade").
				Msg(err.Error())
			return
		}

		client := newClient(conn, ws, payload.ID)

		ws.addClient(client)

		go client.readMessages()
		go client.writeMessages()

	}
}

func (ws *WSChat) addClient(client *client) {
	ws.Lock()
	defer ws.Unlock()
	ws.clients[client] = true
}

func (ws *WSChat) removeClient(client *client) {
	ws.Lock()
	defer ws.Unlock()
	if _, ok := ws.clients[client]; ok {
		client.conn.Close()
		delete(ws.clients, client)
	}
}
