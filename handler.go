package main

import (
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
