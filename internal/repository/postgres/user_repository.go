package postgres

import (
	"database/sql"
	"errors"
	"github.com/firstProject/internal/domain"
	_ "github.com/lib/pq"
	"log/slog"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int64) (*domain.User, error)
	Update(user *domain.User) error
}

type userRepository struct {
	db *sql.DB // Подключение к базе данных
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create добавляет нового пользователя в базу данных
func (r *userRepository) Create(user *domain.User) error {
	const op = "repository.postgres.UserRepository.Create"

	query := `
		INSERT INTO users (first_name, last_name, email)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Email).Scan(&user.ID)
	if err != nil {
		slog.Error(op, "ошибка при создании пользователя", "error", err)
		return err
	}
	slog.Info(op, "пользователь успешно создан", "user_id", user.ID)
	return nil
}

// GetByID получает информацию о пользователе по его ID
func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	const op = "repository.postgres.UserRepository.GetByID"
	user := &domain.User{}

	query := `
		SELECT id, first_name, last_name, email
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info(op, "пользователь не найден", "user_id", id)
			return nil, nil
		}
		slog.Error(op, "ошибка при получении пользователя", "error", err)
		return nil, err
	}

	slog.Info(op, "пользователь успешно получен", "user_id", id)
	return user, nil
}

// Update обновляет информацию о существующем пользователе
func (r *userRepository) Update(user *domain.User) error {
	const op = "repository.postgres.UserRepository.Update"

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3
		WHERE id = $4`

	result, err := r.db.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		slog.Error(op, "ошибка при обновлении пользователя", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error(op, "ошибка при получении количества обновленных строк", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Info(op, "пользователь для обновления не найден", "user_id", user.ID)
		return sql.ErrNoRows
	}

	slog.Info(op, "пользователь успешно обновлен", "user_id", user.ID)
	return nil
}
