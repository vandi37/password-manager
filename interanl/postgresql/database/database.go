package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorOpeningDataBase = "error opining database"
	CheckingConnection   = "checking database connection failed"
)

type DB struct {
	db *sql.DB
}

func New(host string, username string, password string, port int, name string) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("%s://%s:%s@db:%d/%s?sslmode=disable", host, username, password, port, name))
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningDataBase, err, vanerrors.EmptyHandler)
	}
	err = db.Ping()
	if err != nil {
		return nil, vanerrors.NewWrap(CheckingConnection, err, vanerrors.EmptyHandler)
	}
	return &DB{db: db}, nil
}