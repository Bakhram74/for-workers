package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/logger"

	"strings"

	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker jwt.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			logger.Log.Error().
				Str("authMiddleware", "authorizationHeader==0").
				Msg(err.Error())
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			logger.Log.Error().
				Str("authMiddleware", "authorizationHeader len<2").
				Msg(err.Error())
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			logger.Log.Error().
				Str("authMiddleware", "not authorizationTypeBearer").
				Msg(err.Error())
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		if payload.Role == "" {
			httpServer.ErrorResponse(ctx, http.StatusUnauthorized, jwt.ErrInvalidToken.Error())
			logger.Log.Error().
				Str("authMiddleware", "payload.Role").
				Msg(jwt.ErrInvalidToken.Error())
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
func getUserId(ctx *gin.Context) (string, error) {
	payload, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		return "", errors.New("user id not found")
	}
	tokenPayload, ok := payload.(*jwt.Payload)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}
	return tokenPayload.ID, nil
}
