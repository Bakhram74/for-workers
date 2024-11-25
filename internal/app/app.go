package app

import (
	"fmt"

	"github.com/ShamilKhal/shgo/config"

	db "github.com/ShamilKhal/shgo/db/sqlc"
	"github.com/ShamilKhal/shgo/internal/controller/http"
	"github.com/ShamilKhal/shgo/internal/service"
	"github.com/ShamilKhal/shgo/pkg/client/postgres"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/ShamilKhal/shgo/pkg/httpServer"
	"github.com/ShamilKhal/shgo/pkg/jwt"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(config *config.Config) {

	jwtMaker, err := jwt.NewJWTMaker(config.Token.TokenSecretKey)
	if err != nil {
		logger.Log.Fatal().Msgf("Cannot create token maker: %s", err.Error())
	}

	logger.Log.Info().Msg("Postgresql initializing")
	pool := initDB(config)
	defer pool.Close()

	err = RunMigration(config)
	if err != nil {
		logger.Log.Fatal().Msgf("Migration error: %s", err.Error())
	}

	store := db.NewStore(pool)

	logger.Log.Info().Msg("Redis initializing")

	redisClient, err := redis.InitRedis(config)
	if err != nil {
		logger.Log.Fatal().Msg(err.Error())
	}
	defer redis.Close()
	redis.CreateFetchChatBetweenIndex()

	service := service.NewService(service.Deps{
		Store:    store,
		JwtMaker: jwtMaker,
		Config:   config,
		Redis:    redisClient,
	})

	handler := http.NewHandler(service, jwtMaker, config).Init()

	logger.Log.Info().Msgf("Runnig app server at %s", config.HTTP.Address)
	
	srv := httpServer.NewServer(config, logger.Middleware(handler))
	if err := srv.Run(); err != nil {
		logger.Log.Error().Msgf("Error occurred while running http server: %s\n", err.Error())
	}

}

func initDB(cfg *config.Config) *pgxpool.Pool {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database,
		cfg.Storage.SSLMode)

	conn, err := postgres.New(dsn)
	if err != nil {
		logger.Log.Fatal().Msgf("Failed to initialize db: %s", err.Error())
	}

	return conn
}
