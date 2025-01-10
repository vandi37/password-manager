// TODO: finish service
package service

import (
	"context"

	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/repo/repo"
	"github.com/vandi37/password-manager/pkg/password"
)

type Service struct {
	repo            *repo.Repo
	passwordService *password.PasswordService
}

func New(repo *repo.Repo, passwordService *password.PasswordService) *Service {
	return &Service{repo: repo, passwordService: passwordService}
}

func (s *Service) Decrypt(master, cipherText, nonce []byte) ([]byte, error) {
	return s.passwordService.Decrypt(master, cipherText, nonce)
}

func (s *Service) UpdatePasswordUsername(ctx context.Context, password_id int, username string) error {
	return s.repo.PasswordRepo.UpdateUsername(ctx, password_id, username)
}

func (s *Service) NewUser(ctx context.Context, id int64, password string) error {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return err
	}

	return s.repo.UserRepo.Create(ctx, id, hash)
}

func (s *Service) CheckUserPassword(ctx context.Context, id int64, password string) (bool, error) {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return false, err
	}

	return s.repo.UserRepo.Compare(ctx, hash, id)
}

func (s *Service) UpdateUser(ctx context.Context, id int64, password string) error {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return err
	}

	return s.repo.UserRepo.Update(ctx, id, hash)
}

func (s *Service) RemoveUser(ctx context.Context, id int64) error {
	return s.repo.UserRepo.Delete(ctx, id)
}

func (s *Service) UserExists(ctx context.Context, id int64) (bool, error) {
	return s.repo.UserRepo.Exist(ctx, id)
}

func (s *Service) NewPassword(ctx context.Context, id int64, master string, password string, company string, username string) error {
	res, nonce, err := s.passwordService.Encrypt([]byte(master), []byte(password))
	if err != nil {
		return err
	}
	return s.repo.PasswordRepo.Create(ctx, module.Password{
		Company:  company,
		Username: username,
		Password: res,
		Nonce:    nonce,
		UserId:   id,
	})
}

func (s *Service) GetPasswordByUserId(ctx context.Context, id int64) ([]module.Password, error) {
	return s.repo.PasswordRepo.GetByUserId(ctx, id)
}
