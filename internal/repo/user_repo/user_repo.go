package user_repo

import (
	"context"
	"github.com/vandi37/password-manager/internal/postgresql/module"

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

func (r *UserRepo) Create(ctx context.Context, user module.User) error {
	stmt, err := r.db.PrepareContext(ctx, `insert into users (id, password, key, nonce) values ($1, $2, $3, $4);`)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.Id, user.Password, user.Key, user.Nonce)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *UserRepo) Update(ctx context.Context, user module.User) error {
	stmt, err := r.db.PrepareContext(ctx, `update users set password = $1, key = $2, nonce = $3 where id = $4;`)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.Password, user.Key, user.Nonce, user.Id)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *UserRepo) Compare(ctx context.Context, password []byte, id int64) ([]byte, []byte, bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select password = $1, key, nonce from users where id = $2`)
	if err != nil {
		return nil, nil, false, vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, password, id)
	if err != nil {
		return nil, nil, false, vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	defer rows.Close()

	rows.Next()

	var ok bool
	var key []byte
	var nonce []byte

	err = rows.Scan(&ok, &key, &nonce)
	if err != nil {
		return nil, nil, false, vanerrors.Wrap(repo.ErrorScanning, err)
	}

	return key, nonce, ok, nil
}

func (r *UserRepo) Get(ctx context.Context, id int64) (*module.User, error) {
	stmt, err := r.db.PrepareContext(ctx, `select password, key, nonce from users where id = $1`)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorPreparing, err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorExecuting, err)
	}
	defer rows.Close()
	rows.Next()

	var res module.User
	err = rows.Scan(&res.Password, &res.Key, &res.Nonce)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorScanning, err)
	}
	res.Id = id
	return &res, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *UserRepo) Exist(ctx context.Context, id int64) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select coalesce( (select 1 from users where id = $1), 0 );`)
	if err != nil {
		return false, vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return false, vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	defer rows.Close()

	rows.Next()

	var res bool

	err = rows.Scan(&res)
	if err != nil {
		return false, vanerrors.Wrap(repo.ErrorScanning, err)
	}

	return res, nil
}
