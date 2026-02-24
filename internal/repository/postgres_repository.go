package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/models"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	var id int

	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (user_name, email, password, access_level) VALUES ($1, $2, $3, $4) RETURNING id",
		user.UserName, user.Email, user.Password, user.AccessLevel,
	).Scan(&id)

	return id, err
}
