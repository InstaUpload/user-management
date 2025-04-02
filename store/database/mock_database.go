package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/InstaUpload/user-management/types"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Setup(dbConfig *types.DatabaseConfig) error {
	// Do migrations in the database.
	m, err := migrate.New(
		fmt.Sprintf("file://%s", dbConfig.MigrationsFolder),
		dbConfig.GetConnectionString())
	if err != nil {
		log.Fatalf("Error in migrating %v", err)
		return err
	}
	if err := m.Up(); err != nil {
		log.Fatalf("Migration error, %v", err)
	}
	migratedVersion, dirty, err := m.Version()
	if err != nil {
		log.Fatalf("Error in migrating %v", err)
	}
	log.Printf("Migrated to version %d (dirty=%t)", migratedVersion, dirty)
	return nil
}

// Function to run when test are over.
func KillPostgresContainer(container *tcpg.PostgresContainer) {
	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func CreatePostgresContainer(ctx context.Context, dbConfig *types.DatabaseConfig) (*tcpg.PostgresContainer, error) {
	// Create a new postgres test container.
	dbName := dbConfig.Name
	dbUser := dbConfig.User
	dbPassword := dbConfig.Password
	postgresContainer, err := tcpg.Run(ctx,
		"postgres:16-alpine",
		tcpg.WithDatabase(dbName),
		tcpg.WithUsername(dbUser),
		tcpg.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}
	return postgresContainer, nil
}
