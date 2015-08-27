package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

var queries = map[string]string{
	"user-by-email-and-active": "select * from users where email = $1 and active = $2 limit 1",
}

func Init(pgConnection string) error {
	var err error

	db, err := sqlx.Open("postgres", pgConnection)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(8)

	DB = db
	return nil
}

func Get(obj interface{}, queryKey string, args ...interface{}) error {
	err := DB.Get(obj, queries[queryKey], args...)
	if err != nil {
		return err
	}
	return nil
}
