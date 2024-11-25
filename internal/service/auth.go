package service

import (
	"context"
	"errors"

	"github.com/ShamilKhal/shgo/config"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService struct {
	store    db.Store
	config   *config.Config
	redis    *redis.Redis
	jwtMaker jwt.Maker
}

func NewAuthService(store db.Store, config *config.Config, redis *redis.Redis, jwtMaker jwt.Maker) *AuthService {
	return &AuthService{
		store:    store,
		config:   config,
		redis:    redis,
		jwtMaker: jwtMaker,
	}
}

type stickerData struct {
	Sticker string `json:"sticker"`
}

func (s *AuthService) Login(ctx context.Context, phone string) (*stickerData, error) {

	err := redis.SetLimit(ctx, phone)
	if err != nil {
		return nil, err
	}

	user, err := s.store.GetUserByPhone(ctx, phone)
	if err == nil {
		return &stickerData{Sticker: user.ID}, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		args := db.CreateUserParams{
			ID:    uuid.NewString(),
			Phone: phone,
		}

		newUser, err := s.store.CreateUser(ctx, args)
		if err == nil {
			return &stickerData{Sticker: newUser.ID}, nil
		}
	}

	return nil, err

}

type tokenData struct {
	AccessToken  string
	RefreshToken string
	User         db.User
}

func (s *AuthService) AuthCreateUser(ctx context.Context, userID, name, imageUrl, statusText string) (tokenData, error) {

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return tokenData{}, err
	}

	if user.Role != "guest" {
		return tokenData{}, errors.New("user already exists")
	}

	user, err = s.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:         user.ID,
		Role:       pgtype.Text{Valid: true, String: "user"},
		Name:       pgtype.Text{Valid: true, String: name},
		ImageUrl:   pgtype.Text{Valid: true, String: imageUrl},
		StatusText: pgtype.Text{Valid: true, String: statusText},
	})
	if err != nil {
		return tokenData{}, err
	}
	tokens, err := s.createToken(user)
	if err != nil {
		return tokenData{}, err
	}

	data := tokenData{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}
	return data, nil
}

func (s *AuthService) AuthVerifyUser(ctx context.Context, userID string, pincode string) (tokenData, error) {
	value, err := redis.GetPin(ctx, userID)
	if err != nil {
		return tokenData{}, err
	}
	if value.Pincode != pincode {
		return tokenData{}, errors.New("redis: wrong pincode")
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return tokenData{}, err
	}

	if err := redis.DeletePinValue(ctx, userID); err != nil {
		return tokenData{}, err
	}

	tokens, err := s.createToken(user)
	if err != nil {
		return tokenData{}, err
	}

	return tokenData{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) AuthRefreshToken(ctx context.Context, refreshToken string) (authToken, error) {
	payload, err := s.jwtMaker.VerifyToken(refreshToken)
	if err != nil {
		return authToken{}, err
	}
	if payload.Role != "" {
		return authToken{}, jwt.ErrInvalidToken
	}
	user, err := s.store.GetUserByID(ctx, payload.ID)
	if err != nil {
		return authToken{}, err
	}
	tokens, err := s.createToken(user)
	if err != nil {
		return authToken{}, err
	}

	return authToken{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

type authToken struct {
	AccessToken  string
	RefreshToken string
	Payload      *jwt.Payload
}

func (s *AuthService) createToken(user db.User) (authToken, error) {

	accessToken, payload, err := s.jwtMaker.CreateToken(
		user,
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		return authToken{}, err
	}

	userId := db.User{
		ID: user.ID,
	}
	refreshToken, _, err := s.jwtMaker.CreateToken(
		userId,
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		return authToken{}, err
	}

	return authToken{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Payload:      payload,
	}, nil
}
