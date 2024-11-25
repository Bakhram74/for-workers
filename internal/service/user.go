package service

import (
	"context"
	"errors"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	store db.Store
	redis *redis.Redis
}

func NewUserService(store db.Store, redis *redis.Redis) *UserService {
	return &UserService{
		store: store,
		redis: redis,
	}
}

func (s *UserService) UpdateUserPhone(ctx context.Context, userID string, pincode string) (db.User, error) {
	value, err := redis.GetPin(ctx, userID)
	if err != nil {
		return db.User{}, err
	}
	if value.Pincode != pincode {
		return db.User{}, errors.New("redis: wrong pincode")
	}

	if err := redis.DeletePinValue(ctx, userID); err != nil {
		return db.User{}, err
	}

	arg := db.UpdateUserParams{
		Phone: pgtype.Text{Valid: true, String: value.Phone},
		ID:    userID,
	}
	return s.store.UpdateUser(ctx, arg)
}

func (s *UserService) UpdateUserImg(ctx context.Context, userID string, imgUrl string) (db.User, error) {

	arg := db.UpdateUserParams{
		ImageUrl: pgtype.Text{Valid: true, String: imgUrl},
		ID:       userID,
	}
	return s.store.UpdateUser(ctx, arg)
}

func (s *UserService) UpdateUserData(ctx context.Context, userID string, name string, statusText string) (db.User, error) {

	var arg db.UpdateUserParams

	switch {
	case statusText == "" && name != "":
		{
			arg = db.UpdateUserParams{
				Name: pgtype.Text{Valid: true, String: name},
				ID:   userID,
			}
		}
	case name == "" && statusText != "":
		{
			arg = db.UpdateUserParams{
				StatusText: pgtype.Text{Valid: true, String: statusText},
				ID:         userID,
			}
		}
	default:
		{
			arg = db.UpdateUserParams{
				Name:       pgtype.Text{Valid: true, String: name},
				StatusText: pgtype.Text{Valid: true, String: statusText},
				ID:         userID,
			}
		}
	}
	return s.store.UpdateUser(ctx, arg)
}
