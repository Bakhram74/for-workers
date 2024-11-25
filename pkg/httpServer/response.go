package httpServer

import "github.com/gin-gonic/gin"

type Response struct {
	Data any `json:"data,omitempty"`
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, Error{
		Err{
			Message: msg,
		},
	})
}

type Error struct {
	Err `json:"error"`
}

type Err struct {
	Message string `json:"message"`
}
