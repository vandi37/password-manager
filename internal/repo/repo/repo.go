package repo

import (
	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/repo/password_repo"
	"github.com/vandi37/password-manager/internal/repo/user_repo"
)

type Repo struct {
	UserRepo     *user_repo.UserRepo
	PasswordRepo *password_repo.PasswordRepo
}

func New(db *database.DB) *Repo {
	return &Repo{
		UserRepo:     user_repo.New(db),
		PasswordRepo: password_repo.New(db),
	}
}
