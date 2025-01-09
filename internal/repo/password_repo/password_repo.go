package password_repo

import (
	"context"
	"database/sql"

	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/repo"
	"github.com/vandi37/vanerrors"
)

type PasswordRepo struct {
	db *database.DB
}

func New(db *database.DB) *PasswordRepo {
	return &PasswordRepo{db: db}
}

func (r *PasswordRepo) Create(ctx context.Context, password module.Password) error {
	stmt, err := r.db.PrepareContext(ctx, `insert into passwords (company, username, password, nonce, user_id) values ($1, $2, $3, $4, $5);`)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorPreparing, err, vanerrors.EmptyHandler)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, password.Company, password.Username, password.Password, password.Nonce, password.UserId)
	if err != nil {
		return vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	return nil
}

func (r *PasswordRepo) GetByUserId(ctx context.Context, id int64) ([]module.Password, error) {
	rows, err := r.db.QueryContext(ctx, "select company, username, password, nonce, user_id from passwords where user_id = $1", id)
	if err != nil {
		return nil, vanerrors.NewWrap(repo.ErrorExecuting, err, vanerrors.EmptyHandler)
	}

	defer rows.Close()

	return scanPasswordRows(rows)
}

func scanPasswordRows(rows *sql.Rows) ([]module.Password, error) {
	res := []module.Password{}
	for rows.Next() {
		var password module.Password
		err := rows.Scan(&password.Company, &password.Username, &password.Password, &password.Nonce, &password.UserId)
		if err != nil {
			return nil, vanerrors.NewWrap(repo.ErrorScanning, err, vanerrors.EmptyHandler)
		}
		res = append(res, password)
	}
	return res, nil
}
