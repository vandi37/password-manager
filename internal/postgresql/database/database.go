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

type DB struct {
	*sql.DB
}

func New(ctx context.Context, username string, password string, host string, port int, name string) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, name))
	if err != nil {
		return nil, vanerrors.Wrap(ErrorOpeningDataBase, err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, vanerrors.Wrap(CheckingConnection, err)
	}
	return &DB{DB: db}, nil
}

func (db *DB) Close(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		err = db.DB.Close()
		cancel()
	}()

	<-ctx.Done()
	return
}

func (db *DB) Init(ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
	create table if not exists users (
		id bigint not null primary key,
		password bytea not null,
		key bytea not null,
		nonce bytea not null,
		created timestamp with time zone not null default current_timestamp
	);
	
	create table if  not exists passwords (
		id serial not null primary key,
		company text not null,
		username text not null,
		password bytea not null,
		nonce bytea not null,
		user_id bigint not null references users(id) on delete cascade,
		created timestamp with time zone not null default current_timestamp
	);
	`)

	if err != nil {
		return vanerrors.Wrap(ErrorCreateTable, err)
	}

	return nil
}
