package db

import (
	"log"
	"os"

	"testing"

	"github.com/ShamilKhal/shgo/pkg/client/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/shupir?sslmode=disable"
)

func initDB() *pgxpool.Pool {
	conn, err := postgres.New(dbSource)
	if err != nil {
		log.Fatal("failed to initialize db: ")
	}
	return conn
}

var testQueries *Queries

func TestMain(m *testing.M) {
	conn := initDB()
	testQueries = New(conn)

	os.Exit(m.Run())
}
