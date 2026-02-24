package repository

import (
	"context"
	"github.com/tjasha/Recipe_manager/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
}
