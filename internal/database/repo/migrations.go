package repo

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (pgRepo *PgRepo) ApplyMigrations() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}
	if err := goose.Up(pgRepo.db.Master, "migrations"); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}
	return nil
}
