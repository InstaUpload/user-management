package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/broker"
	"github.com/InstaUpload/user-management/service"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/store/database"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello from User Management System")
	// Loading environment variables.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Setting up the database configration.
	dbConfig := types.DatabaseConfig{
		User:             utils.GetEnvString("DATABASEUSER", "user"),
		Password:         utils.GetEnvString("DATABASEPASSWORD", "user"),
		Name:             utils.GetEnvString("DATABASENAME", "user"),
		MaxOpenConns:     utils.GetEnvInt("DATABASEOPENCONNS", 5),
		MaxIdleConns:     utils.GetEnvInt("DATABASEIDLECONNS", 5),
		MaxIdleTime:      utils.GetEnvString("DATABASEIDLETIME", "1m"),
		MigrationsFolder: "./migrations",
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Name)
	log.Println(connectionString)
	dbConfig.SetConnectionString(connectionString)
	db, err := database.New(&dbConfig)
	if err != nil {
		log.Fatalf("Can not connect to database %v \n", err)
	}
	defer db.Close()
	database.Setup(&dbConfig)
	// Setting up kafka producer.
	brokers := []string{"localhost:9092"}
	producer, err := broker.ConnectProducer(brokers)
	if err != nil {
		log.Fatalf("Can not connect to kafka producer %v \n", err)
	}
	defer producer.Close()
	// Setting up store.
	dbStore := store.NewStore(db)
	// Setting up service.
	grpcService := service.NewService(&dbStore)
	// Setting up broker sender.
	sender := broker.NewSender(producer)
	// Setting up handler.
	handler := NewHandler(&grpcService, &sender)
	// Setting up grpc server.
	grpcAddress := utils.GetEnvString("USERSERVICEADDRESS", "localhost:5003")
	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("Error in listing for tcp connection")
	}
	defer l.Close()
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, handler)

	log.Printf("gRPC server for user management microservice running in port %s", grpcAddress)
	if err := s.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
