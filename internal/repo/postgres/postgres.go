package postgres

import (
	"github.com/wb-go/wbf/dbpg"
)

type PgRepo struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *PgRepo {
	return &PgRepo{
		db: db,
	}
}
