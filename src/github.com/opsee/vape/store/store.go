package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB
)

func Init(pgConnection string) error {
	var err error

	db, err = sqlx.Open("postgres", pgConnection)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(8)

	return nil
}
