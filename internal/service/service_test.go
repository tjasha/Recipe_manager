package service

import (
	"context"
	"testing"

	"github.com/tjasha/Recipe_manager/internal/models"
)

type MockRepo struct{}

func (m *MockRepo) CreateUser(ctx context.Context, user *models.User) (int, error) {
	return 3, nil
}

func TestCreateUser(t *testing.T) {
	mockRepo := &MockRepo{}
	service := NewService(mockRepo)

	id, err := service.CreateUser(context.Background(), "testUser", "user@gmail.com", "pass", 1)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	if id != 3 {
		t.Errorf("Expected user ID 3, got %d", id)
	}
}
