package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/firstProject/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserRepository struct {
	users map[int64]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[int64]*domain.User),
	}
}

func (m *mockUserRepository) Create(user *domain.User) error {
	user.ID = int64(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(id int64) (*domain.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *mockUserRepository) Update(user *domain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return sql.ErrNoRows
	}
	m.users[user.ID] = user
	return nil
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	handler := NewUserHandler(mockRepo)

	user := domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	rw := httptest.NewRecorder()

	handler.CreateUser(rw, req)

	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rw.Code)
	}

	var response domain.User
	json.NewDecoder(rw.Body).Decode(&response)

	if response.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	handler := NewUserHandler(mockRepo)

	user := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	mockRepo.Create(user)

	req := httptest.NewRequest("GET", "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()

	handler.GetUser(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	var response domain.User
	json.NewDecoder(rw.Body).Decode(&response)

	if response.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", response.ID)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	handler := NewUserHandler(mockRepo)

	user := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	mockRepo.Create(user)

	updatedUser := domain.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
	}

	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()

	handler.UpdateUser(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	var response domain.User
	json.NewDecoder(rw.Body).Decode(&response)

	if response.FirstName != "Jane" {
		t.Errorf("Expected updated first name 'Jane', got '%s'", response.FirstName)
	}
}