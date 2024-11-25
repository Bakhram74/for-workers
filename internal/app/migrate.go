package app

import (
	"errors"
	"fmt"

	"time"

	"github.com/ShamilKhal/shgo/config"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func RunMigration(cfg *config.Config) error {
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database,
		cfg.Storage.SSLMode)

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://db/migration", databaseURL)
		if err == nil {
			break
		}

		logger.Log.Debug().Msgf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return errors.New("migrate: postgres connect error: %s")
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Log.Info().Msg("Migrate: no change")
		return nil
	}

	if err != nil {
		return tryCloseMigrateWithError(m, err)
	}

	logger.Log.Info().Msg("Migrate: up success")

	return tryCloseMigrateWithError(m, nil)
}

func tryCloseMigrateWithError(m *migrate.Migrate, err error) error {
	var resultErr error
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		resultErr = fmt.Errorf("failed to close source, err: %w", sourceErr)
	}
	if databaseErr != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("failed to close database, err: %w", databaseErr))
	}
	return errors.Join(err, resultErr)
}
