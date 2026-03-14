package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// CreateUser creates new user and returns its ID.
// It's gonna be used only for Goole login, it doesn't save password.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (user_name, email, google_id, access_level, created_at, modified_at) 
		 VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		 RETURNING id`,
		user.UserName, user.Email, user.GoogleID, user.AccessLevel,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByEmail retrieves a user from the database by their email address.
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_name, email, google_id, access_level, created_at, modified_at, state 
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.UserName, &user.Email, &user.GoogleID, &user.AccessLevel, &user.CreatedAt, &user.ModifiedAt, &user.State)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) AuthenticateUser(ctx context.Context, user *model.User) (int, string, int, error) {
	var id, accessLevel int
	var UserName string

	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (user_name, email, password, access_level) VALUES ($1, $2, $3, $4) RETURNING id, user_name, access_level",
		user.UserName, user.Email, user.Password, user.AccessLevel,
	).Scan(&id, &UserName, &accessLevel)

	return id, UserName, accessLevel, err
}
