package main

import (
	"context"
	"log"

	pb "github.com/InstaUpload/common/api"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
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
	// Sending message to send welcome email.
	var sendWelcome types.SendWelcomeEmailKM
	sendWelcome.Name = in.Name
	sendWelcome.Email = in.Email
	if err := h.messageSender.Email.SendWelcome(&sendWelcome); err != nil {
		log.Printf("Error sending welcome email: %s", err.Error())
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
	token, err := h.grpcService.User.Login(ctx, &userData)
	// if error found return empity string and error.
	if err != nil {
		return nil, err
	}
	// else return token and empty error.
	return &pb.LoginUserResponse{Token: token}, nil
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
		Role:       userData.Role.Name,
		CreatedOn:  utils.TimeToString(&userData.CreatedOn),
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

func (h *Handler) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	email := in.Email
	_, err := h.grpcService.User.ResetPassword(ctx, email)
	if err != nil {
		return nil, err
	}
	return &pb.ResetUserPasswordResponse{Msg: "Email sent to registered email address."}, nil
}

func (h *Handler) UpdateUserPassword(ctx context.Context, in *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordResponse, error) {
	token := in.Token
	password := in.Password

	if err := h.grpcService.User.UpdatePassword(ctx, token, password); err != nil {
		return nil, err
	}
	return &pb.UpdateUserPasswordResponse{Msg: "Password updated."}, nil
}

func (h *Handler) VerifyUser(ctx context.Context, in *pb.VerifyUserRequest) (*pb.VerifyUserResponse, error) {
	token := in.Token
	if err := h.grpcService.User.Verify(ctx, token); err != nil {
		return nil, err
	}
	return &pb.VerifyUserResponse{Msg: "User verified."}, nil
}

func (h *Handler) SendVerification(ctx context.Context, in *pb.SendVerificationUserRequest) (*pb.SendVerificationUserResponse, error) {
	// TODO: Maybe add the Kafka producer to send email in handler struct.
	// TODO: get token from send verification menthod and send it to the user.
	// TODO: using kafka.
	token, err := h.grpcService.User.SendVerification(ctx)
	if err != nil {
		return nil, err
	}
	var sendVerification types.SendVerificationKM
	sendVerification.Token = token
	if err := h.messageSender.Email.SendVerification(&sendVerification); err != nil {
		return nil, err
	}
	return &pb.SendVerificationUserResponse{}, nil
}

func (h *Handler) AddEditorUser(ctx context.Context, in *pb.AddEditorUserRequest) (*pb.AddEditorUserResponse, error) {
	userId := in.UserId
	if err := h.grpcService.User.AddEditor(ctx, userId); err != nil {
		return nil, err
	}
	return &pb.AddEditorUserResponse{}, nil
}
