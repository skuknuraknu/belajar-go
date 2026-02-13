package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration parameters for the application.
// It includes settings for the server, environment, and database.
type Config struct {
	// Server configuration
	Port int    // The port on which the server will listen.
	Env  string // The environment in which the application is running (e.g., "development", "production").

	// Database configuration
	DBDriver string // The database driver (e.g., "postgres", "mysql").
	DBDSN    string // The Data Source Name (DSN) for connecting to the database.
}

// Load loads the application configuration from command-line flags.
// It returns a Config struct populated with the provided values or defaults.
func Load() (*Config, error) {
	var cfg Config

	// Server flags
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	// Database flags
	flag.StringVar(&cfg.DBDriver, "db-driver", "mysql", "Database driver (e.g., postgres, mysql)")
	// Note: In a real application, the DSN might be loaded from an environment variable
	// for security, but as per requirements, we'll use a flag for now.
	// A default DSN is provided for convenience.
	defaultDBDSN := "greenlight:pa55word@tcp(localhost:3306)/greenlight?parseTime=true"
	flag.StringVar(&cfg.DBDSN, "db-dsn", defaultDBDSN, "Database DSN (Data Source Name)")

	// Parse the command-line flags.
	flag.Parse()

	// Basic validation
	if cfg.Env != "development" && cfg.Env != "staging" && cfg.Env != "production" {
		return nil, fmt.Errorf("invalid environment: %s. Must be one of: development, staging, production", cfg.Env)
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d. Must be between 1 and 65535", cfg.Port)
	}

	// Allow overriding DSN via an environment variable for better security practice.
	// This is a common pattern: flags for defaults, env vars for overrides.
	if dbDSNEnv := os.Getenv("GREENLIGHT_DB_DSN"); dbDSNEnv != "" {
		cfg.DBDSN = dbDSNEnv
	}

	// Allow overriding port via an environment variable.
	if portEnv := os.Getenv("GREENLIGHT_PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid GREENLIGHT_PORT environment variable: %w", err)
		}
		cfg.Port = port
	}

	return &cfg, nil
}
