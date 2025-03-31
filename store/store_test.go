package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/InstaUpload/user-management/store/database"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
	"github.com/joho/godotenv"
)

func PSting() {
	log.Printf("Store package")
}

var MockStore Store

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
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
	if err = db.Ping(); err != nil {
		log.Fatal("Can not connect to database")
	}
	database.Setup(&testDbConfig)
	MockStore = NewStore(db)
	exitCode := m.Run()
	database.KillPostgresContainer(container)
	os.Exit(exitCode)

}
