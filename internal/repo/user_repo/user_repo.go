package user_repo

import (
	"context"

	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/repo"
	"github.com/vandi37/vanerrors"
)

type UserRepo struct {
	db *database.DB
}

func New(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, id int64, password []byte) error {
	stmt, err := r.db.PrepareContext(ctx, `insert into users (id, password) values ($1, $2);`)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, password)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, id int64, password []byte) error {
	stmt, err := r.db.PrepareContext(ctx, `update users set password = $1 where id = $2;`)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, password, id)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return nil
}

func (r *UserRepo) Compare(ctx context.Context, password []byte, id int64) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select password = $1 from users where id = $2`)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, password, id)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	rows.Next()

	var res bool

	err = rows.Scan(&res)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
	}

	return res, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return nil
}

func (r *UserRepo) Exist(ctx context.Context, id int64) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select coalesce( (select 1 from users where id = $1), 0 );`)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	rows.Next()

	var res bool

	err = rows.Scan(&res)
	if err != nil {
		return false, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
	}

	return res, nil
}
