package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/loadlab-go/usersvc/model"
)

var db *DB

func mustInitDB(dn, dsn string) {
	db = &DB{sqlx.MustOpen(dn, dsn)}
	db.dbx.SetMaxOpenConns(10)
	logger.Info("DB initialized")
}

type DB struct {
	dbx *sqlx.DB
}

func (d *DB) CreateUser(name, password string) (*model.User, error) {
	u := &model.User{
		Name:     name,
		Password: password,
	}
	err := d.dbx.QueryRow(`INSERT INTO "user" (name, password) VALUES ($1, $2) RETURNING id`, name, password).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}
