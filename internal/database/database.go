package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	TimeZone        string
	Debug           bool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func ConfigFromEnv() Config {
	return Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		Name:            getEnv("DB_NAME", "payment"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		TimeZone:        getEnv("DB_TIMEZONE", "Asia/Singapore"),
		Debug:           getEnvBool("DB_DEBUG", true),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}
}

func Open(config Config) (*gorm.DB, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
		config.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newQueryLogger(config.Debug),
	})
	if err != nil {
		return nil, fmt.Errorf("membuka database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("mengakses connection pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("mengakses connection pool: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("menutup connection pool: %w", err)
	}

	return nil
}

func validateConfig(config Config) error {
	if config.MaxOpenConns < 1 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS harus lebih besar dari 0")
	}
	if config.MaxIdleConns < 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS tidak boleh negatif")
	}
	if config.MaxIdleConns > config.MaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS tidak boleh melebihi DB_MAX_OPEN_CONNS")
	}
	if config.ConnMaxLifetime < 0 {
		return fmt.Errorf("DB_CONN_MAX_LIFETIME tidak boleh negatif")
	}
	if config.ConnMaxIdleTime < 0 {
		return fmt.Errorf("DB_CONN_MAX_IDLE_TIME tidak boleh negatif")
	}

	return nil
}

func newQueryLogger(debug bool) logger.Interface {
	logLevel := logger.Warn
	if debug {
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "[GORM] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
