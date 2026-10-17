package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"backend/internal/auth"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Register(
	ctx context.Context, input RegisterInput,
) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	API_key, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, err
	}
	return s.repository.Create(
		ctx,
		input.Username,
		input.Email,
		string(hash),
		API_key,
	)
}

func (s *Service) Login(
	ctx context.Context,
	input LoginInput,
) (*User, error) {
	user, err := s.repository.Get(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
