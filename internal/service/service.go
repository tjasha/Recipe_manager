package service

import (
	"context"

	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateUser(ctx context.Context, userName, email, password string, accessLevel int) (int, error) {
	user := &model.User{
		UserName:    userName,
		Email:       email,
		Password:    password,
		AccessLevel: accessLevel,
	}

	return s.repo.CreateUser(ctx, user)
}
