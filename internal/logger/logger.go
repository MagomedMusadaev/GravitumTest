package logger

import (
	"fmt"
	"github.com/firstProject/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"strings"
)

// NewLogger создает новый логгер с настроенной ротацией и уровнем логирования
func NewLogger(cfg *config.LogConfig) *slog.Logger {
	// Настройка ротации логов
	logWriter := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // максимальный размер в мегабайтах
		MaxBackups: cfg.MaxBackups, // максимальное количество файлов
		MaxAge:     cfg.MaxAge,     // максимальный возраст в днях
		Compress:   cfg.Compress,   // сжимать ротированные файлы
	}

	// Определение уровня логирования
	var level slog.Level
	switch strings.ToUpper(cfg.Level) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Настройка обработчика логов
	var handler slog.Handler

	// В режиме разработки используем текстовый формат
	if cfg.Environment == "development" {
		handler = slog.NewTextHandler(logWriter, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// В продакшене используем JSON формат
		handler = slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}

// SetupGlobalLogger устанавливает глобальный логгер
func SetupGlobalLogger(cfg *config.LogConfig) error {
	logger := NewLogger(cfg)
	slog.SetDefault(logger)

	// Проверка возможности записи в файл
	if err := os.MkdirAll(cfg.GetLogDir(), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	return nil
}