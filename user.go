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

func (h *Handler) UpdateUserRole(ctx context.Context, in *pb.UpdateUserRoleRequest) (*pb.UpdateUserRoleResponse, error) {
	role := in.RoleName
	userId := in.UserId
	if err := h.grpcService.User.UpdateRole(ctx, userId, role); err != nil {
		return nil, err
	}
	return &pb.UpdateUserRoleResponse{Msg: "User role updated"}, nil
}

func (h *Handler) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordRespons, error) {
	email := in.Email
	_, err := h.grpcService.User.ResetPassword(ctx, email)
	if err != nil {
		return nil, err
	}
	return &pb.ResetUserPasswordResponse{Msg: "Email send to registered email address."}, nil
}

func (h *Handler) UpdateUserPassword(ctx context.Context, in *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordResponse, error) {
	token := in.Token
	password := in.Password

	if err := h.grpcService.User.UpdatePassword(ctx, email, password); err != nil {
		return nil, err
	}
	return &pb.UpdateUserPasswordResponse{Msg: "Password updated."}, nil
}
