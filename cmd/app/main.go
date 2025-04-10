package main

import (
	"database/sql"
	"fmt"
	"github.com/firstProject/internal/config"
	"github.com/firstProject/internal/handler"
	"github.com/firstProject/internal/logger"
	"github.com/firstProject/internal/repository/postgres"
	"github.com/firstProject/internal/routes"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	const op = "cmd.app.main"

	// Загрузка конфигурации
	cfg := config.NewConfig()

	// Настройка логирования
	if err := logger.SetupGlobalLogger(&cfg.Log); err != nil {
		fmt.Printf("Не удалось настроить logger: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Запуск приложения", "environment", cfg.Env)

	// Установка соединения с базой данных Psql
	db, err := sql.Open("postgres", cfg.DB.GetConnectionString())
	if err != nil {
		slog.Error(op, "Ошибка подключения к базе данных:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err = db.Ping(); err != nil {
		slog.Error(op, "Ошибка при проверке соединения с базой данных:", err)
		os.Exit(1)
	}

	// Инициализация репозитория и создание обработчика HTTP-запросов
	userRepo := postgres.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Настройка маршрутизации HTTP-запросов
	r := mux.NewRouter()
	routes.SetupUserRoutes(r, userHandler)

	// Запуск HTTP-сервера на указанном порту
	fmt.Printf("Сервер запущен и прослушивает порт %s\n", cfg.Server.Port)
	if err = http.ListenAndServe(cfg.Server.Port, r); err != nil {
		slog.Error(op, "Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
}
