package v1

import (
	"net/http"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (r *Routes) initChatRoutes(api *gin.RouterGroup) {
	chat := api.Group("/chat").Use(authMiddleware(r.jwtMaker))
	{
		chat.GET("/contact-list", r.contactList)
		chat.GET("/chat-history", r.chatHistory)
	}

}

// @Summary Contact lists
// @Security ApiKeyAuth
// @Description Contact lists
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} []entity.ContactList
// @Failure      401,500  {object}  httpServer.Error
// @Router /chat/contact-list [get]
func (r *Routes) contactList(ctx *gin.Context) {
	userID, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("contactList", "getUserId(ctx)").
			Msg(err.Error())
		return
	}

	contactList, err := r.service.IChat.GetContactList(userID)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("contactList", "IChat.GetContactList").
			Msg(err.Error())
		return
	}

	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: contactList,
	})
}

// @Summary Chat history
// @Security ApiKeyAuth
// @Description Chat history
// @Tags chat
// @Accept json
// @Produce json
// @Param			user	query		string		true	"used id"
// @Success 200 {object} []entity.Chat
// @Failure      401,400,500  {object}  httpServer.Error
// @Router /chat/chat-history [get]
func (r *Routes) chatHistory(ctx *gin.Context) {
	userID1, err := getUserId(ctx)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		logger.Log.Error().
			Str("chatHistory", "getUserId(ctx)").
			Msg(err.Error())
		return
	}
	userID2 := ctx.Query("user")
	if userID2 == "" {
		httpServer.ErrorResponse(ctx, http.StatusBadRequest, "query param user is missing")
		logger.Log.Error().
			Str("chatHistory", "ctx.Query(user)").
			Msg("query param user is missing")
		return
	}

	// chat between timerange fromTS toTS
	// where TS is timestamp
	// 0 to positive infinity
	fromTS, toTS := "0", "+inf"

	chat, err := r.service.IChat.GetChatHistory(userID1, userID2, fromTS, toTS)
	if err != nil {
		httpServer.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		logger.Log.Error().
			Str("chatHistory", "IChat.GetChatHistory").
			Str("userID1", userID1).
			Str("userID2", userID2).
			Msg(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, httpServer.Response{
		Data: chat,
	})
}
