package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Security SecurityConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Env  string
}

type SecurityConfig struct {
	SessionSecret   string
	SessionDuration time.Duration
}

func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Security.SessionSecret == "change-this-to-a-random-secret-in-production" && c.Server.Env == "production" {
		return fmt.Errorf("SESSION_SECRET must be changed in production")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found. Using environmental variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			DBName:   getEnv("DB_NAME", "gobank"),
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "dxtrr"),
			Password: getEnv("DB_PASSWORD", ""),
			Port:     getEnv("DB_PORT", "5432"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Security: SecurityConfig{
			SessionSecret:   getEnv("SESSION_SECRET", "change-this-to-a-random-secret-in-production"),
			SessionDuration: getDurationEnv("SESSION_DURATION", 24) * time.Hour,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// method over Config to get Data source Name
func (c *Config) GetDNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env != "production"
}
func (c *Config) IsProduction() bool {
	return c.Server.Env != "development"
}

//Helpers

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
func getIntEnv(key string, defualtValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defualtValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defualtValue
	}
	return value
}

func getDurationEnv(key string, defaultValue int) time.Duration {

	return time.Duration(getIntEnv(key, defaultValue))
}
