package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(database *sql.DB) error {
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create db driver: %s", err)
	}
	s, err := WithInstance(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to open embedded migrations: %s", err)
	}
	m, err := migrate.NewWithInstance("embed", s, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrations: %s", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %s", err)
	}
	return nil
}
