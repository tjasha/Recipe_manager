package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tjasha/Recipe_manager/internal/config"
)

//NOV PREDLOG
//// V testni datoteki:
//type MockRepository struct {}
//
//func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
//	// Ni klica v bazo!
//	// Samo vrnemo vnaprej določen rezultat za ta test.
//	return 123, nil
//}
//
//// ... v testni funkciji ...
//mockRepo := &MockRepository{}
//app := &handler.Application{
//DB: mockRepo, // Uporabimo lažni repozitorij!
//// ... ostale odvisnosti
//}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	config := config.LoadConfig()

	dbURL := config.TEST_DB_URL
	if dbURL == "" {
		t.Fatal("TEST_DB_URL environment variable not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}

	// Clear the users table before each test
	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		pool.Close()
		t.Fatalf("Failed to truncate users table: %v", err)
	}
	return pool
}

func teardownTestDB(t *testing.T, pool *pgxpool.Pool) {
	// Clear the users table after each test
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Errorf("Failed to truncate users table during teardown: %v", err)
	}
	pool.Close()
}

//
//func TestCreateUser(t *testing.T) {
//	pool := setupTestDB(t)
//	defer teardownTestDB(t, pool)
//
//	repo := NewPostgresRepository(pool)
//
//	id, googleId, accessLevel, err := repo.AuthenticateUser(context.Background(), &model.User{
//		UserName: "testUser",
//		Email:    "user@test.com",
//		//Password:    "pass",
//		AccessLevel: 1,
//		GoogleID:    "testID",
//	})
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if id == 0 {
//		t.Fatal("expected valid ID")
//	}
//}
