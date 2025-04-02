package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/store/database"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Hello from User Management System Business Test.")
	testDbConfig := types.DatabaseConfig{
		User:             utils.GetEnvString("TESTDATABASEUSER", "tester"),
		Password:         utils.GetEnvString("TESTDATABASEPASSWORD", "user"),
		Name:             utils.GetEnvString("TESTDATABASENAME", "userdb"),
		MaxOpenConns:     utils.GetEnvInt("TESTDATABASEOPENCONNS", 5),
		MaxIdleConns:     utils.GetEnvInt("TESTDATABASEIDLECONNS", 5),
		MaxIdleTime:      utils.GetEnvString("TESTDATABASEIDLETIME", "1m"),
		MigrationsFolder: "../migrations",
	}
	ctx := context.Background()
	container, err := database.CreatePostgresContainer(ctx, &testDbConfig)
	if err != nil {
		log.Fatalf("Can not create postgres container")
		return
	}
	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Error in getting connection string: %v", err)
		return
	}
	log.Printf("connection string %s", connectionString)
	testDbConfig.SetConnectionString(connectionString)
	db, err := database.New(&testDbConfig)
	if err != nil {
		log.Fatalf("Can not create new database")
	}
	database.Setup(&testDbConfig)
	store.MockStore = store.NewStore(db)
	validate = validator.New()
	exitCode := m.Run()
	database.KillPostgresContainer(container)
	os.Exit(exitCode)

}
