package ws

import (
	"encoding/json"
	"log"

	"time"

	"github.com/ShamilKhal/shgo/internal/entity"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second //TODO refactor
	pingInterval = (pongWait * 9) / 10
)

type clientList map[*client]bool

type client struct {
	conn   *websocket.Conn
	ws     *WSChat
	userID string
	egress chan entity.Chat
}

func newClient(conn *websocket.Conn, ws *WSChat, userID string) *client {
	return &client{
		conn:   conn,
		ws:     ws,
		userID: userID,
		egress: make(chan entity.Chat),
	}
}

func (client *client) pongHandler(pongMsg string) error {
	return client.conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (client *client) readMessages() {
	defer func() {
		client.ws.removeClient(client)
	}()

	client.conn.SetReadLimit(512)

	if err := client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Log.Error().
			Str("func", "readMessages()").
			Msg(err.Error())
		return
	}

	client.conn.SetPongHandler(client.pongHandler)

	for {
		_, payload, err := client.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Error().
					Str("func", "readMessages()").
					Msgf("error reading message: %s", err.Error())
			}
			break
		}

		chat := entity.Chat{}

		if err := json.Unmarshal(payload, &chat); err != nil {
			logger.Log.Error().
				Str("func", "readMessages()").
				Msgf("error marshalling message: %s", err.Error())
			break
		}
		chat.From = client.userID
		chat.Timestamp = time.Now().Unix()

		id, err := client.ws.service.IChat.CreateChat(chat)
		if err != nil {
			logger.Log.Error().
				Str("func", "readMessages()").
				Msgf("error while saving chat in redis %s", err.Error())
			return
		}
		chat.ID = id
		for c := range client.ws.clients {
			if c.userID == chat.From || c.userID == chat.To {
				c.egress <- chat
			}
		}

	}
}

func (client *client) writeMessages() {

	ticker := time.NewTicker(pingInterval)

	defer func() {
		client.ws.removeClient(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					logger.Log.Error().
						Str("func", "writeMessages()").
						Msgf("connection closed: %s", err.Error())
				}
				return
			}
			if err := client.conn.WriteJSON(message); err != nil {
				log.Println(err)
				logger.Log.Error().
					Str("func", "writeMessages()").
					Msg(err.Error())
			}
		case <-ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.Log.Error().
					Str("func", "writeMessages()").
					Msgf("writemsg: %s", err.Error())
				return
			}
		}

	}
}
