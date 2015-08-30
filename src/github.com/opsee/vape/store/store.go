package store

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Tx struct {
	*sqlx.Tx
}

var DB *sqlx.DB

var queries = map[string]string{
	// users
	"user-by-email-and-active": "select * from users where email = $1 and active = $2",
	"user-by-id":               "select * from users where id = $1",
	"delete-user-by-id":        "delete from users where id = $1",
	"update-user":              "update users set name = :name, email = :email, password_hash = :password_hash where id = :id",
	"insert-user":              "insert into users (org_id, email, name, verified, active, password_hash) values (:org_id, :email, :name, :verified, :active, :password_hash) returning *",

	// signups
	"signup-by-id":  "select * from signups where id = $1",
	"insert-signup": "insert into signups (email, name) values (:email, :name) returning *",
	"list-signups":  "select * from signups limit $1 offset $2",
	"claim-signup":  "update signups set claimed = true where id = $1",

	// orgs
	"insert-new-org": "insert into orgs (name) values (NULL) returning id",

	// bastions
	"insert-bastion":           "insert into bastions (password_hash, org_id) values (:password_hash, :org_id) returning *",
	"bastion-by-id-and-active": "select * from bastions where id = $1 and active = $2",
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
	return DB.Get(obj, queries[queryKey], args...)
}

func Select(obj interface{}, queryKey string, args ...interface{}) error {
	return DB.Select(obj, queries[queryKey], args...)
}

func Exec(queryKey string, args ...interface{}) (sql.Result, error) {
	return DB.Exec(queries[queryKey], args...)
}

func NamedExec(queryKey string, arg interface{}) (sql.Result, error) {
	return DB.NamedExec(queries[queryKey], arg)
}

func NamedQuery(queryKey string, arg interface{}) (*sqlx.Rows, error) {
	return DB.NamedQuery(queries[queryKey], arg)
}

func NamedInsert(queryKey string, arg interface{}) error {
	rows, err := NamedQuery(queryKey, arg)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.StructScan(arg); err != nil {
			return err
		}
	}
	return nil
}

func Beginx() (*Tx, error) {
	tx, err := DB.Beginx()
	return &Tx{tx}, err
}

func (tx *Tx) Get(obj interface{}, queryKey string, args ...interface{}) error {
	return tx.Tx.Get(obj, queries[queryKey], args...)
}

func (tx *Tx) Select(obj interface{}, queryKey string, args ...interface{}) error {
	return tx.Tx.Select(obj, queries[queryKey], args...)
}

func (tx *Tx) Exec(queryKey string, args ...interface{}) (sql.Result, error) {
	return tx.Tx.Exec(queries[queryKey], args...)
}

func (tx *Tx) NamedExec(queryKey string, arg interface{}) (sql.Result, error) {
	return tx.Tx.NamedExec(queries[queryKey], arg)
}

func (tx *Tx) NamedQuery(queryKey string, arg interface{}) (*sqlx.Rows, error) {
	return tx.Tx.NamedQuery(queries[queryKey], arg)
}

func (tx *Tx) NamedInsert(queryKey string, arg interface{}) error {
	rows, err := tx.NamedQuery(queryKey, arg)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err = rows.StructScan(arg); err != nil {
			return err
		}
	}
	return nil
}
