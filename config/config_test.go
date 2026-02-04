package config

import (
	// "os"
	"testing"
)

// func TestLoad(t *testing.T) {
// 	os.Setenv("DB_PASSWORD", "testpassword")
// 	os.Setenv("SESSION_SECRET", "testsecret")

// 	cfg, err := Load()

// 	if err != nil {
// 		t.Fatalf("Expected no error, but got %v ", err)
// 	}

// 	if cfg.Database.Password != "testpassword" {
// 		t.Errorf("Expected DB_PASSWORD to be 'testpassword' but got %s", cfg.Database.Password)
// 	}

// 	if cfg.Security.SessionSecret != "testsecret" {
// 		t.Errorf("Expected SESSION_SECRET to be 'testsecret' but got %s", cfg.Security.SessionSecret)
// 	}
// 	if cfg.Server.Port != "8080" {
// 		t.Errorf("Expected default PORT to be '8080', got '%s'", cfg.Server.Port)
// 	}

// 	if cfg.Database.Host != "localhost" {
// 		t.Errorf("Expected default DB_HOST to be 'localhost', got '%s'", cfg.Database.Host)
// 	}
// 	os.Unsetenv("DB_PASSWORD")
// 	os.Unsetenv("SESSION_SECRET")
// }

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			name:      "valid config",
			shouldErr: false,
			config: &Config{
				Database: DatabaseConfig{
					Password: "password",
					DBName:   "testdb",
				},
				Security: SecurityConfig{
					SessionSecret: "some-secret",
				},
			},
		},
		{
			name:      "missing password",
			shouldErr: true,
			config: &Config{
				Database: DatabaseConfig{
					Password: "",
					DBName:   "testdb",
				},
				Security: SecurityConfig{
					SessionSecret: "some-secret",
				},
			},
		},
		{
			name:      "missing database name",
			shouldErr: true,
			config: &Config{
				Database: DatabaseConfig{
					Password: "password",
					DBName:   "",
				},
				Security: SecurityConfig{
					SessionSecret: "some-secret",
				},
			},
		},
		{
			name:      "insecure session secret in production",
			shouldErr: true,
			config: &Config{
				Database: DatabaseConfig{
					Password: "password",
					DBName:   "testdb",
				},
				Server: ServerConfig{
					Env: "production",
				},
				Security: SecurityConfig{
					SessionSecret: "change-this-to-a-random-secret-in-production",
				},
			},
		},
		{
			name:      "insecure session secret in development",
			shouldErr: false,
			config: &Config{
				Database: DatabaseConfig{
					Password: "password",
					DBName:   "testdb",
				},
				Server: ServerConfig{
					Env: "development",
				},
				Security: SecurityConfig{
					SessionSecret: "change-this-to-a-random-secret-in-production",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if err == nil && tt.shouldErr {
				t.Errorf("expected errors but got none")
			}
			if err != nil && !tt.shouldErr {
				t.Errorf("expected no errors but got %v", err)
			}
		})
	}
}
