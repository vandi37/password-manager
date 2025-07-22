package service

import (
	"context"
	"github.com/vandi37/password-manager/pkg/generate"
	"github.com/vandi37/password-manager/pkg/logger"
	"go.uber.org/zap"

	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/repo/repo"
	"github.com/vandi37/password-manager/pkg/password"
)

const keyLength = 32

const Namespace = "service-namespace"

var Field = zap.Dict(Namespace)

type Service struct {
	repo            *repo.Repo
	passwordService *password.Service
}

func New(repo *repo.Repo, passwordService *password.Service) *Service {
	return &Service{repo: repo, passwordService: passwordService}
}

func (s *Service) Decrypt(master, cipherText, nonce []byte) ([]byte, error) {
	return s.passwordService.Decrypt(master, cipherText, nonce)
}

func (s *Service) UpdatePassword(ctx context.Context, passwordId int, password string, key string) error {
	res, nonce, err := s.passwordService.Encrypt([]byte(key), []byte(password))
	if err != nil {
		logger.Debug(ctx, "Failed to encrypt password", zap.Error(err), zap.Int(logger.PasswordId, passwordId), Field)
		return err
	}
	err = s.repo.PasswordRepo.Update(ctx, passwordId, res, nonce)
	if err != nil {
		logger.Debug(ctx, "Failed to update password", zap.Error(err), zap.Int(logger.PasswordId, passwordId), Field)
		return err
	}
	logger.Debug(ctx, "Updated password", zap.Int(logger.PasswordId, passwordId), Field)
	return nil
}

func (s *Service) UpdatePasswordUsername(ctx context.Context, passwordId int, username string) error {
	err := s.repo.PasswordRepo.UpdateUsername(ctx, passwordId, username)
	if err != nil {
		logger.Debug(ctx, "Failed to update password username", zap.Error(err), zap.Int(logger.PasswordId, passwordId), zap.String(logger.Username, username), Field)
		return err
	}
	logger.Debug(ctx, "Updated password username", zap.Int(logger.PasswordId, passwordId), zap.String(logger.Username, username), Field)
	return nil
}

func (s *Service) NewUser(ctx context.Context, id int64, password string) error {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		logger.Debug(ctx, "Failed to hash password", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	key, nonce, err := s.passwordService.Encrypt([]byte(password), []byte(generate.Password(keyLength, true, true, true, true)))
	if err != nil {
		logger.Debug(ctx, "Failed to encrypt key", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}

	err = s.repo.UserRepo.Create(ctx, module.User{
		Id:       id,
		Password: hash,
		Key:      key,
		Nonce:    nonce,
	})
	if err != nil {
		logger.Debug(ctx, "Failed to create user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}

	logger.Debug(ctx, "Created user", zap.Int64(logger.UserId, id), Field)
	return nil
}

func (s *Service) CheckUserPassword(ctx context.Context, id int64, password string) (string, bool, error) {
	hash, err := s.passwordService.Hash(password)
	if err != nil {
		logger.Debug(ctx, "Failed to hash password", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return "", false, err
	}

	key, nonce, ok, err := s.repo.UserRepo.Compare(ctx, hash, id)
	if err != nil {
		logger.Debug(ctx, "Failed to compare user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return "", false, err
	}
	res, err := s.Decrypt([]byte(password), key, nonce)
	if err != nil {
		logger.Debug(ctx, "Failed to decrypt key", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return "", false, err
	}
	logger.Debug(ctx, "Checked user password", zap.Int64(logger.UserId, id), zap.Bool(logger.IsOk, ok), Field)
	return string(res), ok, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, password string, lastPassword string) error {
	user, err := s.repo.UserRepo.Get(ctx, id)
	if err != nil {
		logger.Debug(ctx, "Failed to get user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}

	res, err := s.Decrypt([]byte(lastPassword), user.Key, user.Nonce)
	if err != nil {
		logger.Debug(ctx, "Failed to decrypt key", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	key, nonce, err := s.passwordService.Encrypt([]byte(password), res)
	if err != nil {
		logger.Debug(ctx, "Failed to encrypt key", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}

	hash, err := s.passwordService.Hash(password)
	if err != nil {
		logger.Debug(ctx, "Failed to hash password", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}

	err = s.repo.UserRepo.Update(ctx, module.User{
		Id:       id,
		Password: hash,
		Key:      key,
		Nonce:    nonce,
	})
	if err != nil {
		logger.Debug(ctx, "Failed to update user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	logger.Debug(ctx, "Updated user", zap.Int64(logger.UserId, id), Field)
	return nil
}

func (s *Service) RemovePassword(ctx context.Context, passwordId int) error {
	err := s.repo.PasswordRepo.Remove(ctx, passwordId)
	if err != nil {
		logger.Debug(ctx, "Failed to remove password", zap.Error(err), zap.Int(logger.PasswordId, passwordId), Field)
		return err
	}
	logger.Debug(ctx, "Removed password", zap.Int(logger.PasswordId, passwordId), Field)
	return nil
}

func (s *Service) RemoveUser(ctx context.Context, id int64) error {
	err := s.repo.UserRepo.Delete(ctx, id)
	if err != nil {
		logger.Debug(ctx, "Failed to remove user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	logger.Debug(ctx, "Removed user", zap.Int64(logger.UserId, id), Field)
	return nil
}

func (s *Service) UserExists(ctx context.Context, id int64) (bool, error) {
	ok, err := s.repo.UserRepo.Exist(ctx, id)
	if err != nil {
		logger.Debug(ctx, "Failed to check user", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return false, err
	}
	logger.Debug(ctx, "User exists", zap.Int64(logger.UserId, id), zap.Bool(logger.IsOk, ok), Field)
	return ok, nil
}

func (s *Service) NewPassword(ctx context.Context, id int64, key string, password string, company string, username string) error {
	res, nonce, err := s.passwordService.Encrypt([]byte(key), []byte(password))
	if err != nil {
		logger.Debug(ctx, "Failed to encrypt password", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	err = s.repo.PasswordRepo.Create(ctx, module.Password{
		Company:  company,
		Username: username,
		Password: res,
		Nonce:    nonce,
		UserId:   id,
	})
	if err != nil {
		logger.Debug(ctx, "Failed to create password", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return err
	}
	logger.Debug(ctx, "Created password", zap.Int64(logger.UserId, id), Field)
	return nil
}

func (s *Service) GetPasswordsByUserId(ctx context.Context, id int64) ([]module.Password, error) {
	passwords, err := s.repo.PasswordRepo.GetByUserId(ctx, id)
	if err != nil {
		logger.Debug(ctx, "Failed to get passwords", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return nil, err
	}
	logger.Debug(ctx, "Get passwords", zap.Int64(logger.UserId, id), Field)
	return passwords, nil
}

func (s *Service) GetPasswordsByCompany(ctx context.Context, id int64, company string) ([]module.Password, error) {
	passwords, err := s.repo.PasswordRepo.GetByCompany(ctx, id, company)
	if err != nil {
		logger.Debug(ctx, "Failed to get passwords", zap.Error(err), zap.Int64(logger.UserId, id), Field)
		return nil, err
	}
	logger.Debug(ctx, "Get passwords", zap.Int64(logger.UserId, id), Field)
	return passwords, nil
}
