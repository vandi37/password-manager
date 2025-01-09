// TODO: finish service
package service

import (
	"context"

	"github.com/vandi37/password-manager/internal/repo/user_repo"
	"github.com/vandi37/password-manager/pkg/password"
	//"github.com/vandi37/password-manager/pkg/logger"
)

type Service struct {
	userRepo        *user_repo.UserRepo
	passwordService *password.PasswordService
	// logger *logger.Logger
}

func New(userRepo *user_repo.UserRepo, passwordService *password.PasswordService) *Service {
	return &Service{userRepo: userRepo, passwordService: passwordService}
}

func (s *Service) NewUser(ctx context.Context, id int64, password string) error {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return err
	}

	return s.userRepo.Create(ctx, id, hash)
}

func (s *Service) CheckUserPassword(ctx context.Context, id int64, password string) (bool, error) {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return false, err
	}

	return s.userRepo.Compare(ctx, hash, id)
}

func (s *Service) UpdateUser(ctx context.Context, id int64, password string) error {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		return err
	}

	return s.userRepo.Update(ctx, id, hash)
}

func (s *Service) RemoveUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *Service) UserExists(ctx context.Context, id int64) (bool, error) {
	return s.userRepo.Exist(ctx, id)
}
