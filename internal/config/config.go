package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	ServerPort     string
	PublicKey      string
	PrivateKey     string
	GoogleClientId string
	TokenDuration  time.Duration
}

var AppConfig Config

func Load() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}

	AppConfig = Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", ""),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", ""),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		ServerPort:     getEnv("SERVER_PORT", "3000"),
		PublicKey:      getEnv("PUBLIC_KEY", ""),
		PrivateKey:     getEnv("PRIVATE_KEY", ""),
		GoogleClientId: getEnv("GOOGLE_CLIENT_ID", ""),
		TokenDuration:  time.Duration(getEnvAsInt("TOKEN_DURATION_MINUTES", 60*24*30)) * time.Minute,
	}

	log.Info("Configuration loaded successfully")
	return nil
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
