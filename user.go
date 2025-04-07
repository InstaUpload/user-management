package main

import (
	"context"

	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/types"
)

func (h *Handler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// convert in to UserPayload.
	userData := types.CreateUserPayload{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
	if err := h.grpcService.User.Create(ctx, &userData); err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{Msg: "User created successfully"}, nil
}

func (h *Handler) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	// convert in to LoginUserPayload.
	userData := types.LoginUserPayload{
		Email:    in.Email,
		Password: in.Password,
	}
	// call grpcService.User.Login() and check for err.
	token, err = h.grpcService.User.Login(ctx, &userData)
	// if error found return empity string and error.
	if err != nil {
		return nil, err
	}
	// else return token and empty error.
	return &pb.LoginUserResponse{token: token}, nil
}

func (h *Handler) AuthUser(ctx context.Context, in *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	token := in.Token
	userData, err := h.grpcService.User.Auth(ctx, token)
	if err != nil {
		return nil, err
	}
	userRes := pb.AuthUserResponse{
		Id:         userData.Id,
		Name:       userData.Name,
		Email:      userData.Email,
		CreatedAt:  userData.CreatedAt,
		IsVerified: userData.IsVerified,
	}
	return &userRes, nil
}
