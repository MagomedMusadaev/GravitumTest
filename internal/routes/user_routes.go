package routes

import (
	"github.com/firstProject/internal/handler"
	"github.com/gorilla/mux"
)

// SetupUserRoutes настраивает маршруты для пользовательских эндпоинтов
func SetupUserRoutes(r *mux.Router, userHandler handler.UserHandler) {
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
}
