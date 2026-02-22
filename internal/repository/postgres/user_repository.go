package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/models"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (user_name, email, password, access_level) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	return r.pool.QueryRow(ctx, query,
		user.UserName,
		user.Email,
		user.Password,
		user.AccessLevel,
	).Scan(&user.ID)
}
