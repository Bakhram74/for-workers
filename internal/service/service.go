package service

import (
	"context"

	"github.com/ShamilKhal/shgo/config"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/ShamilKhal/shgo/internal/entity"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/ShamilKhal/shgo/pkg/jwt"
)

type IAuthorization interface {
	Login(ctx context.Context, phone string) (*stickerData, error)
	AuthVerifyUser(ctx context.Context, sticker string, pincode string) (tokenData, error)
	AuthCreateUser(ctx context.Context, userID, name, imageUrl, statusText string) (tokenData, error)
	AuthRefreshToken(ctx context.Context, refreshToken string) (authToken, error)
}

type IUser interface {
	UpdateUserPhone(ctx context.Context, userID string, pincode string) (db.User, error)
	UpdateUserImg(ctx context.Context, userID string, imgUrl string) (db.User, error)
	UpdateUserData(ctx context.Context, userID string, name string, statusText string) (db.User, error)
}
type IVehicle interface {
	CreateVehicle(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error)
	DeleteVehicle(ctx context.Context, arg db.DeleteVehicleParams) error
	UpdateVehicle(ctx context.Context, arg db.UpdateVehicleParams) (db.Vehicle, error)
	GetVehicleByID(ctx context.Context, arg db.GetVehicleByIDParams) (db.Vehicle, error)
	FindVehicle(ctx context.Context, number string, region int, limit, offset int) ([]db.Vehicle, int, error)
}

type IChat interface {
	GetContactList(id string) ([]entity.ContactList, error)
	GetChatHistory(userID1, userID2, fromTS, toTS string) ([]entity.Chat, error)
	CreateChat(c entity.Chat) (string, error)
}

type Service struct {
	IAuthorization
	IUser
	IVehicle
	IChat
}

func NewService(deps Deps) *Service {
	authService := NewAuthService(deps.Store, deps.Config, deps.Redis, deps.JwtMaker)
	userService := NewUserService(deps.Store, deps.Redis)
	vehicleService := NewVehicleService(deps.Store)
	chatServise := NewChatService(deps.Redis)
	return &Service{
		IAuthorization: authService,
		IUser:          userService,
		IVehicle:       vehicleService,
		IChat:          chatServise,
	}

}

type Deps struct {
	Store    db.Store
	JwtMaker jwt.Maker
	Config   *config.Config
	Redis    *redis.Redis
}
