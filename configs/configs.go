package configs

import (
	"fmt"
	"os"
)

// Environment holds all configuration for the service
type Environment struct {
	AppPort      string
	AppEnv       string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	DBSSLMode    string
	JWTSecret    string
	JWTExpiresIn string
}

// SetEnv sets the global environment config
func SetEnv(env Environment) {
	sharedEnv = env
}

// GetEnv returns the global environment config
func GetEnv() Environment {
	return sharedEnv
}

var sharedEnv Environment

// LoadConfigs loads configuration from environment variables
func LoadConfigs() Environment {
	env := Environment{
		AppPort:      envOrDefault("APP_PORT", "8080"),
		AppEnv:       envOrDefault("APP_ENV", "development"),
		DBHost:       envOrDefault("DB_HOST", "localhost"),
		DBPort:       envOrDefault("DB_PORT", "5432"),
		DBName:       envOrDefault("DB_NAME", "task_management"),
		DBUser:       envOrDefault("DB_USER", "postgres"),
		DBPassword:   envOrDefault("DB_PASSWORD", "postgres"),
		DBSSLMode:    envOrDefault("DB_SSLMODE", "disable"),
		JWTSecret:    envOrDefault("JWT_SECRET", "Ycahyu2WEIy0lRKOO7n6kFT67KCK46KNccVWbBYnP7g"),
		JWTExpiresIn: envOrDefault("JWT_EXPIRES_IN", "24h"),
	}
	SetEnv(env)
	return env
}

// DSN returns the database connection string
func (c Environment) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode)
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
