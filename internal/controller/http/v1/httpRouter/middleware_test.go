package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"
	"time"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/ShamilKhal/shgo/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker jwt.Maker,
	authorizationType string,
	id string,
	name string,
	phone string,
	role string,
	duration time.Duration,
) {
	user := db.User{
		ID:    id,
		Phone: phone,
		Role:  role,
		Name:  name,
	}
	token, payload, err := tokenMaker.CreateToken(user, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	phone := utils.RandomNumbers(11)
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker jwt.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, id, "Alex", phone, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", id, phone, "Alex", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedJWTType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, id, phone, "Alex", "", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
				addAuthorization(t, request, tokenMaker, "", id, phone, "Alex", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker jwt.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, id, phone, "Alex", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			router := gin.New()
			authPath := "/auth"
			tokenMaker, _ := jwt.NewJWTMaker(utils.RandomString(32))
			router.GET(
				authPath,
				authMiddleware(tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, tokenMaker)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
