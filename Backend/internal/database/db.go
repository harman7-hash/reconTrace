package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" 
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib" 
)

// Connect creates and verifies a PostgreSQL connection pool for your main application logic.
func Connect(ctx context.Context, databaseURI string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(databaseURI)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}

	dbConfig.MaxConns = 10
	dbConfig.MinConns = 5
	dbConfig.MaxConnIdleTime = time.Hour
	dbConfig.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// RunMigrations pushes database schema changes forward right before the pool starts up.
func RunMigrations(databaseURI string) error {
	// 1. Parse the incoming connection URI string
	config, err := pgx.ParseConfig(databaseURI)
	if err != nil {
		return fmt.Errorf("failed to parse database string config for migrations: %w", err)
	}

	// 2. Open a temporary standard sql.DB connection utilizing the pgx engine underneath
	db := stdlib.OpenDB(*config)
	defer db.Close() // Automatically cleans up this network connection when migrations finish

	// 3. Attach the database instance to the golang-migrate framework driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver instance: %w", err)
	}

	// 4. Point the migration engine to your SQL folder files layout
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not initialize migration engine: %w", err)
	}

	log.Println("Running database migrations via pgx driver...")
	err = m.Up()
	if err != nil {
		// If no modifications have occurred, handle it gracefully without crashing
		if err == migrate.ErrNoChange {
			log.Println("Database schema is already up to date. No changes applied.")
			return nil
		}
		return fmt.Errorf("migration up execution failed: %w", err)
	}

	log.Println("All migrations executed successfully!")
	return nil
}
