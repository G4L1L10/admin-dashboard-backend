package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string

	ServerPort string
	JwtSecret  string
}

var AppConfig *Config

// LoadConfig loads environment variables and initializes AppConfig
func LoadConfig() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found. Using system environment variables...")
	}

	AppConfig = &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "bandroom_admin"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JwtSecret:  getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
