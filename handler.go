package main

import (
	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/broker"
	"github.com/InstaUpload/user-management/service"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	grpcService   *service.Service
	messageSender *broker.Sender
}

func NewHandler(gs *service.Service, ms *broker.Sender) *Handler {
	return &Handler{
		grpcService:   gs,
		messageSender: ms,
	}
}
