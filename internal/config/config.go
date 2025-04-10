package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	DB        DBConfig     // Настройки базы данных
	Server    ServerConfig // Настройки сервера
	Log       LogConfig    // Настройки логирования
	Env       string       // Текущее окружение (development, production, test)
}

// DBConfig содержит параметры подключения к базе данных
type DBConfig struct {
	Host     string // Хост базы данных
	Port     int    // Порт базы данных
	User     string // Имя пользователя
	Password string // Пароль
	DBName   string // Название базы данных
	SSLMode  string // Режим SSL-подключения
}

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Port string // Порт для HTTP-сервера
}

// LogConfig содержит настройки логирования
type LogConfig struct {
	Level      string // Уровень логирования (DEBUG, INFO, WARN, ERROR)
	FilePath   string // Путь к файлу логов
	MaxSize    int    // Максимальный размер файла в мегабайтах
	MaxBackups int    // Максимальное количество файлов бэкапа
	MaxAge     int    // Максимальное время хранения файлов в днях
	Compress   bool   // Сжимать ротированные файлы
	Environment string // Окружение для формата логов
}

// GetLogDir возвращает директорию для логов
func (c *LogConfig) GetLogDir() string {
	return filepath.Dir(c.FilePath)
}

// NewConfig создает новый экземпляр конфигурации с учетом переменных окружения
func NewConfig() *Config {
	config := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "users"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "INFO"),
			FilePath:   getEnv("LOG_FILE", "logs/app.log"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 10),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 5),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
			Compress:   getEnvAsBool("LOG_COMPRESS", true),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Env: getEnv("ENVIRONMENT", "development"),
	}

	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid configuration: %v", err))
	}

	return config
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверка настроек БД
	if c.DB.Host == "" {
		return fmt.Errorf("хост базы данных не может быть пустым")
	}
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("недопустимый порт базы данных")
	}
	if c.DB.User == "" {
		return fmt.Errorf("имя пользователя базы данных не может быть пустым")
	}

	// Проверка настроек сервера
	if c.Server.Port == "" {
		return fmt.Errorf("порт сервера не может быть пустым")
	}

	// Проверка настроек логирования
	validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
	if !validLogLevels[strings.ToUpper(c.Log.Level)] {
		return fmt.Errorf("недопустимый уровень логирования: %s", c.Log.Level)
	}

	// Проверка окружения
	validEnvs := map[string]bool{"development": true, "production": true, "test": true}
	if !validEnvs[c.Env] {
		return fmt.Errorf("недопустимое окружение: %s", c.Env)
	}

	return nil
}

// GetConnectionString формирует строку подключения к базе данных
func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Вспомогательные функции для работы с переменными окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := fmt.Sscanf(value, "%d"); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}
