package service

import (
	"arxivhub/internal/models"
	"arxivhub/internal/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repository: repository,
	}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, params models.RegisterUserRequest) (models.User, error) {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("hashing error: %w", err)
	}

	params.Password = hashedPassword

	return s.repository.CreateUser(ctx, params)
}

func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return s.repository.GetUserByUsername(ctx, username)
}

func (s *UserServiceImpl) LoginUser(ctx context.Context, params models.LoginUserRequest) (models.User, error) {
	user, err := s.GetUserByUsername(ctx, params.Username)
	if err != nil {
		return models.User{}, err
	}

	err = s.CheckPassword(user.HashedPassword, params.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (s *UserServiceImpl) CheckPassword(hashedPassword string, password string) error {
	if err := compare(hashedPassword, password); err != nil {
		return errors.New("incorrect password")
	}

	return nil
}

func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
