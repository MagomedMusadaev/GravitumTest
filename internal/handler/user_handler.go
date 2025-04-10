package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/firstProject/internal/domain"
	"github.com/firstProject/internal/repository/postgres"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

type UserHandler interface {
	CreateUser(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
}

// UserHandler обрабатывает HTTP-запросы для управления пользователями
type userHandler struct {
	userRepo postgres.UserRepository
}

// NewUserHandler создает новый экземпляр обработчика пользователей
func NewUserHandler(userRepo postgres.UserRepository) UserHandler {
	return &userHandler{userRepo: userRepo}
}

// CreateUser обрабатывает POST-запрос на создание нового пользователя
func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "internal/handler/UserHandler.CreateUser"

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		msgErr := "Ошибка при разборе данных пользователя"
		slog.Error(op, msgErr, err)
		http.Error(w, msgErr+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userRepo.Create(&user); err != nil {
		http.Error(w, "Ошибка при создании пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser обрабатывает GET-запрос на получение информации о пользователе по ID
func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "internal/handler/UserHandler.GetUser"

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		msgErr := "Некорректный ID пользователя"
		slog.Error(op, msgErr, err)
		http.Error(w, msgErr, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении данных пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser обрабатывает PUT-запрос на обновление данных пользователя
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "internal/handler/UserHandler.UpdateUser"

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		msgErr := "Некорректный ID пользователя"
		slog.Error(op, msgErr, err)
		http.Error(w, msgErr, http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при разборе данных пользователя: "+err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = userID
	if err = h.userRepo.Update(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при обновлении данных пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
