package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
	"time"
)

var (
	dbConn    *sql.DB
	clearFunc func()
)

func TestMain(m *testing.M) {
	dbName := "arxivhub"
	_, dsn := mustInitPostgresContainer(dbName)

	conn, err := ConnectToDB(dsn)
	if err != nil {
		panic(err)
	}

	migr, err := migrate.New("file://../migration", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		src, db := migr.Close()
		if src != nil {
			panic(src)
		}
		if db != nil {
			panic(db)
		}
	}()

	if err := migr.Up(); err != nil {
		panic(err)
	}

	clearFunc = func() {
		if err := migr.Down(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
		if err := migr.Up(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			}
		}
	}

	dbConn = conn

	os.Exit(m.Run())
}

func mustInitPostgresContainer(name string) (*postgres.PostgresContainer, string) {
	ctx := context.Background()

	dbName := name
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, host, port.Port(), dbName)

	return postgresContainer, dsn
}
