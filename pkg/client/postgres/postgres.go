package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 5
	defaultConnTimeout  = time.Second * 3
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Pool         *pgxpool.Pool
}

func New(url string) (*pgxpool.Pool, error) {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = DoWithAttempts(func() error {
		err := pg.Pool.Ping(context.Background())
		if err != nil {
			log.Printf("Failed to connect to postgres due to error %v... Going to do the next attempt\n", err)
			return err
		}

		return nil
	}, pg.connAttempts, pg.connTimeout)

	if err != nil {
		log.Printf("All attempts are exceeded. Unable to connect to PostgreSQL")
		return nil, err
	}

	return pg.Pool, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}
	return err
}
