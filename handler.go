package main

import (
	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/broker/producer"
	"github.com/InstaUpload/user-management/service"
)

type Handler struct {
	pb.UnimplementedUserServiceServer
	grpcService   *service.Service
	messageSender *producer.Sender
}

func NewHandler(gs *service.Service, ms *producer.Sender) *Handler {
	return &Handler{
		grpcService:   gs,
		messageSender: ms,
	}
}
