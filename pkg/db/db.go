package db

import (
	"context"
	"database/sql"
	"embed"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

var conn *pgx.Conn

// Includes the migrations into the build.
//
//go:embed migrations
var embedMigrations embed.FS

// Initdb initializes the database.
func Initdb() {
	var connErr error
	conn, connErr = pgx.Connect(context.Background(), config.Options.DB.Address)

	// If the database connection was unsuccessful, exit the application with an error.
	if connErr != nil {
		log.Fatal(connErr)
	}
}

// EnsureLatestVersion ensures that the database is at the latest version by running all migrations.
func EnsureLatestVersion() {
	if !config.Options.DB.Migrations {
		log.Warning("Automatic database migrations are disabled.")
		return
	}

	database, dbErr := sql.Open("pgx", config.Options.DB.Address)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// For embedding the migrations in the binary.
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(database, "migrations")
	if err != nil {
		log.Fatal("Failed to apply new migrations", err)
	}
}

func rollbackTx(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Debug("Failed to rollback transaction", err)
	}
}
