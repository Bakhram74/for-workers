package redis

import (
	"context"
	"fmt"

	"github.com/ShamilKhal/shgo/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *Redis

type Redis struct {
	client *redis.Client
	config *config.Config
}

func InitRedis(config *config.Config) (*Redis, error) {
	addr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Redis.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	redisClient = &Redis{
		client: client,
		config: config,
	}

	return redisClient, nil
}

func Close() {
	redisClient.client.Close()
}
