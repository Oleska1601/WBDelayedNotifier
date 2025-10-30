package repo

import (
	"database/sql"
	"errors"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type PgRepo struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *PgRepo {
	return &PgRepo{
		db: db,
	}
}

func (r *PgRepo) rollbackTransaction(tx *sql.Tx) {
	if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
		zlog.Logger.Error().
			Err(txErr).
			Str("message", "tx.Rollback")
	}
}
