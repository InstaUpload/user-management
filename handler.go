package main

import (
	"context"
	"log"

	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/service"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	grpcService *service.Service
}

func NewHandler(grpcService *service.Service) *Handler {
	return &Handler{
		grpcService: grpcService,
	}
}

func (h *Handler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("Create User function called from grpc")
	return &pb.CreateUserResponse{Msg: "gRPC function called."}, nil
}
