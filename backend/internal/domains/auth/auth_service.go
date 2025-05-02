package auth

import (
	"context"
	"errors"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/dto"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

var (
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service interface {
	Register(ctx context.Context, in dto.RegisterInput) (token string, err error)
	Login(ctx context.Context, in dto.LoginInput) (token string, err error)
}

type service struct {
	repo   user.UserRepositoryInterface
	jwtKey []byte
}

func NewService(repo user.UserRepositoryInterface, jwtKey []byte) Service {
	return &service{repo: repo, jwtKey: jwtKey}
}

func (s *service) Register(ctx context.Context, in dto.RegisterInput) (string, error) {
	u, err := s.repo.GetUserByEmail(in.Email)
	if err == nil && u.ID != 0 {
		return "", ErrEmailTaken
	}

	hash, err := passwords.Hash(in.Password)
	if err != nil {
		return "", err
	}

	userModel := user.User{
		Username: in.Username,
		Email:    in.Email,
		Password: hash,
	}
	if err := s.repo.Create(userModel); err != nil {
		return "", err
	}

	return jwt.GenerateToken(in.Email, int64(userModel.ID), s.jwtKey)
}

func (s *service) Login(ctx context.Context, in dto.LoginInput) (string, error) {
	u, err := s.repo.GetUserByEmail(in.Email)
	if err != nil || u.ID == 0 {
		return "", ErrInvalidCredentials
	}

	if !passwords.Match(in.Password, u.Password) {
		return "", ErrInvalidCredentials
	}

	return jwt.GenerateToken(u.Email, int64(u.ID), s.jwtKey)
}
