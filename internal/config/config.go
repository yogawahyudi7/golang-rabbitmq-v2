package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Server   ServerConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type LoggerConfig struct {
	Level       string
	Format      string
	FileEnabled bool
	LogDir      string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "rabbitmq_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange:   getEnv("RABBITMQ_EXCHANGE", "app_exchange"),
			Queue:      getEnv("RABBITMQ_QUEUE", "app_queue"),
			RoutingKey: getEnv("RABBITMQ_ROUTING_KEY", "app.message"),
		},
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Logger: LoggerConfig{
			Level:       getEnv("LOG_LEVEL", "info"),
			Format:      getEnv("LOG_FORMAT", "json"),
			FileEnabled: getEnv("LOG_FILE_ENABLED", "true") == "true",
			LogDir:      getEnv("LOG_DIR", "logs"),
			MaxSize:     getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxBackups:  getEnvAsInt("LOG_MAX_BACKUPS", 3),
			MaxAge:      getEnvAsInt("LOG_MAX_AGE", 30),
			Compress:    getEnv("LOG_COMPRESS", "true") == "true",
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}