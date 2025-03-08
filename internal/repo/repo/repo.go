package repo

import (
	"context"
	"github.com/vandi37/password-manager/internal/postgresql/module"
)

type Repo struct {
	UserRepo     UserRepo
	PasswordRepo PasswordRepo
}

func New(userRepo UserRepo, passwordRepo PasswordRepo) *Repo {
	return &Repo{
		UserRepo:     userRepo,
		PasswordRepo: passwordRepo,
	}
}

type UserRepo interface {
	Create(ctx context.Context, user module.User) error
	Update(ctx context.Context, user module.User) error
	Get(ctx context.Context, id int64) (*module.User, error)
	Compare(ctx context.Context, password []byte, id int64) ([]byte, []byte, bool, error)
	Delete(ctx context.Context, id int64) error
	Exist(ctx context.Context, id int64) (bool, error)
}
type PasswordRepo interface {
	Create(ctx context.Context, password module.Password) error
	UpdateUsername(ctx context.Context, passwordId int, username string) error
	Update(ctx context.Context, passwordId int, password []byte, nonce []byte) error
	Remove(ctx context.Context, passwordId int) error
	GetByUserId(ctx context.Context, id int64) ([]module.Password, error)
	GetByCompany(ctx context.Context, id int64, company string) ([]module.Password, error)
}
