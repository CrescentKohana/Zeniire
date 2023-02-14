package config

import (
	"os"
)

// DBOptions is the model for database configuration.
type DBOptions struct {
	Migrations bool
	Address    string // postgres://username:password@localhost:5432/database_name
}

// dbMigrationsEnabled loads the env for enabling or disabling the DB migrations.
func dbMigrationsEnabled() bool {
	// Disabled only when explicitly set to false.
	return os.Getenv("ZNRE_DB_MIGRATIONS") != "false"
}

// dbName loads and parses the env of the DB name.
func dbName() string {
	value := os.Getenv("ZNRE_DB_NAME")
	if value == "" {
		return "zeniire"
	}
	return value
}

// dbName loads and parses the env of the DB host address.
func dbHost() string {
	value := os.Getenv("ZNRE_DB_HOST")
	if value == "" {
		return "localhost"
	}
	return value
}

// dbName loads and parses the env of the DB host port.
func dbPort() string {
	value := os.Getenv("ZNRE_DB_PORT")
	if value == "" {
		return "5332"
	}
	return value
}

// dbName loads and parses the env of the DB username.
func dbUser() string {
	value := os.Getenv("ZNRE_DB_USER")
	if value == "" {
		return "admin"
	}
	return value
}

// dbName loads and parses the env of the DB password.
func dbPassword() string {
	value := os.Getenv("ZNRE_DB_PW")
	if value == "" {
		return "test"
	}
	return value
}
