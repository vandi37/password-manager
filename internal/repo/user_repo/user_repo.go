package user_repo

import (
	"context"

	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/repo/errors"
	"github.com/vandi37/vanerrors"
)

type UserRepo struct {
	db *database.DB
}

func New(db *database.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Compare(ctx context.Context, password []byte, id int64) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, `select master = $1 where id = $2`)
	if err != nil {
		return false, vanerrors.NewWrap(errors.ErrorPreparing, err, vanerrors.EmptyHandler)

	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, password, id)

	return false, nil
	// TODO : FINISH
}
