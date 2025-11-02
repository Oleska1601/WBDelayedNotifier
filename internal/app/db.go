package app

import (
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/wb-go/wbf/dbpg"
)

func initDB(cfg *config.PostgresConfig) (*dbpg.DB, error) {
	masterDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
	slaveDSNs := []string{}
	options := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
	db, err := dbpg.New(masterDSN, slaveDSNs, options)
	if err != nil {
		return nil, fmt.Errorf("create a new DB instance: %w", err)
	}
	return db, nil

}
