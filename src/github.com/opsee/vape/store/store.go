package store

import (
        "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
        db *sqlx.DB
        vapeKey []byte
)

func Init(pgConnection string, sharedKey []byte) error {
        var err error

	db, err = sqlx.Open("postgres", pgConnection)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(8)

        vapeKey = sharedKey
        return nil
}
