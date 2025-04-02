package main

import (
	"context"

	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/types"
)

func (h *Handler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// convert in to UserPayload.
	userData := types.UserPayload{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
	if err := h.grpcService.User.Create(ctx, &userData); err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{Msg: "User created successfully"}, nil
}
