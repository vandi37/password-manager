package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorOpeningDataBase = "error opining database"
	CheckingConnection   = "checking database connection failed"
	ErrorCreateTable     = "error to create table"
)

const (
	MAX_DATA = 2047
)

type DB struct {
	*sql.DB
}

func New(ctx context.Context, username string, password string, host string, port int, name string) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, name))
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorOpeningDataBase, err, vanerrors.EmptyHandler)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, vanerrors.NewWrap(CheckingConnection, err, vanerrors.EmptyHandler)
	}
	return &DB{DB: db}, nil
}

func (db *DB) Close(ctx context.Context) (err error) {
	go func() {
		err = db.DB.Close()
	}()

	<-ctx.Done()
	return
}

func (db *DB) Init(ctx context.Context) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(`
	CREATE TABLE  IF NOT EXISTS users (
		id BIGINT NOT NULL,
		password VARCHAR(32),
		created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);
	
	CREATE TABLE IF NOT EXISTS passwords  (
		id SERIAL,
		password VARCHAR(%d),
		nonce CHAR(12),
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`, MAX_DATA))

	if err != nil {
		return vanerrors.NewWrap(ErrorCreateTable, err, vanerrors.EmptyHandler)
	}

	return nil
}
