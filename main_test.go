package main

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
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	log.Printf("Main function called")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Hello from User Management System")
	testDbConfig := types.DatabaseConfig{
		User:         utils.GetEnvString("TESTDATABASEUSER", "user"),
		Password:     utils.GetEnvString("TESTDATABASEPASSWORD", "user"),
		Name:         utils.GetEnvString("TESTDATABASENAME", "user"),
		MaxOpenConns: utils.GetEnvInt("TESTDATABASEOPENCONNS", 5),
		MaxIdleConns: utils.GetEnvInt("TESTDATABASEIDLECONNS", 5),
		MaxIdleTime:  utils.GetEnvString("TESTDATABASEIDLETIME", "1m"),
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
	database.Setup(db)
	store.MockStore = store.NewStore(db)
	exitCode := m.Run()
	database.KillPostgresContainer(container)
	os.Exit(exitCode)

}
