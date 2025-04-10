# User Management REST API

REST API сервис на Go для управления пользователями с использованием PostgreSQL.

## Функциональность

- Создание пользователя
- Получение информации о пользователе
- Обновление данных пользователя

## Технологии

- Go 1.21
- PostgreSQL
- Docker & Docker Compose
- Gorilla Mux (HTTP router)
- Структурированное логирование (slog)

## Запуск проекта

### С использованием Docker

1. Убедитесь, что у вас установлены Docker и Docker Compose
2. Клонируйте репозиторий
3. Запустите приложение:
   ```bash
   docker-compose up --build
   ```

Сервис будет доступен по адресу: http://localhost:8080

### Локальный запуск

1. Установите Go 1.21
2. Установите и настройте PostgreSQL
3. Создайте базу данных и примените миграции из файла `migrations/init.sql`
4. Установите зависимости:
   ```bash
   go mod download
   ```
5. Запустите приложение:
   ```bash
   go run main.go
   ```

## API Endpoints

### Создание пользователя

```http
POST /users

Content-Type: application/json

{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "secure_password123"
}
```

Ответ:
```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "secure_password123"
}
```

Валидация полей:
- first_name: обязательное поле
- last_name: обязательное поле
- email: обязательное поле, должен быть валидным email адресом
- password: обязательное поле, минимум 8 символов
```

### Получение пользователя

```http
GET /users/{id}
```

Ответ:
```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "secure_password123"
}
```

Возможные ответы:
- 200 OK: пользователь найден
- 404 Not Found: пользователь не найден
- 400 Bad Request: некорректный ID пользователя

### Обновление пользователя

```http
PUT /users/{id}

Content-Type: application/json

{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "new_password123"
}
```

Ответ:
```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "new_password123"
}
```

Возможные ответы:
- 200 OK: данные пользователя обновлены
- 404 Not Found: пользователь не найден
- 400 Bad Request: некорректные данные или ID пользователя
```

## Структура проекта

```.
├── cmd/
│   └── app/         # Точка входа приложения
│       └── main.go
├── internal/
│   ├── config/      # Конфигурация приложения
│   ├── domain/      # Бизнес-сущности и интерфейсы
│   ├── handler/     # HTTP обработчики
│   ├── repository/  # Реализация хранения данных
│   │   └── postgres/# Реализация PostgreSQL репозитория
│   └── routes/      # Настройка маршрутизации
├── migrations/      # SQL миграции
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Логирование

В проекте используется структурированное логирование с помощью пакета `log/slog`. Логи выводятся в формате JSON и содержат дополнительную информацию о контексте операций:
- Название операции (op)
- Уровень логирования (error/info)
- Дополнительные метаданные