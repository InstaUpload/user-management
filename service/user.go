package service

import (
	"context"
	"log"
	"slices"
	"strings"

	common "github.com/InstaUpload/common/types"
	"github.com/InstaUpload/user-management/store"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
)

type UserService struct {
	dbstore    *store.Store
	jwtService interface {
		GenerateAuthToken(int64) (string, error)
		ParseAuthToken(string) (int64, error)
		GeneratePasswordToken(int64) (string, error)
		ParsePasswordToken(string) (int64, error)
		GenerateVerifyToken(int64) (string, error)
		ParseVerifyToken(string) (int64, error)
		GenerateEditorReqToken(int64) (string, error)
		ParseEditorReqToken(string) (int64, error)
	}
}

func (s *UserService) Create(ctx context.Context, userPayload *types.CreateUserPayload) error {
	// Add validation for userPayload struct.
	if err := validate.Struct(userPayload); err != nil {
		// return a custome error from errors package. for invalide data.
		log.Printf("invalid user data: %v", err)
		return common.ErrIncorrectDataReceived
	}
	// below section needs to be automated.
	var user types.User
	user.Name = userPayload.Name
	user.Email = userPayload.Email
	user.Password.Text = userPayload.Password

	user.Password.HashPassword()
	// Check if email is one of the super user's.
	if slices.Contains(utils.GetSuperUsers(), userPayload.Email) {
		// if so then assign 2 to user.RoleId
		user.RoleId = 2
	} else {
		// else make user.RoleId = 1
		user.RoleId = 1
	}
	if err := s.dbstore.User.Create(ctx, &user); err != nil {
		// Check error, if already exist return common.ErrDataFound error.
		// write a if block that check if "duplicate key" error, return common.ErrDataFound error.
		if strings.Contains(err.Error(), "duplicate key") {
			return common.ErrDataFound
		}
		log.Printf("err: %s", err.Error())
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, userPayload *types.LoginUserPayload) (string, error) {
	if err := validate.Struct(userPayload); err != nil {
		// return a custome error from errors package. for invalide data.
		log.Printf("invalid user data: %v", err)
		return "", common.ErrIncorrectDataReceived
	}
	// convert userPayload to user.
	user := types.User{
		Email: userPayload.Email,
	}
	// call user.GetUserByEmail() method.
	// check if error is returned from Login().
	// if error is not nil, return common.ErrDataNotFound error.
	if err := s.dbstore.User.GetUserByEmail(ctx, &user); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", common.ErrDataNotFound
		}

		log.Printf("err: %s", err.Error())
		return "", err
	}
	// match user.Password.hash with userPayload.Password.
	// If match fails return common.ErrIncorrectDataReceived error.
	user.Password.Text = userPayload.Password
	if err := user.Password.ComparePassword(); err != nil {
		return "", common.ErrIncorrectDataReceived
	}
	// generate token and return it.
	token, err := s.jwtService.GenerateAuthToken(user.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		// TODO: To be updated to internal server error.
		return "", common.ErrDataNotFound
	}
	// return token.
	return token, nil
}

func (s *UserService) Auth(ctx context.Context, token string) (types.User, error) {
	// get userId from token.
	var user types.User
	userId, err := s.jwtService.ParseAuthToken(token)
	if err != nil {
		// check the error message and return error accordingly.
		if strings.Contains(err.Error(), "token is expired") {
			return user, common.ErrIncorrectDataReceived
		} else if strings.Contains(err.Error(), "signature is invalid") {
			return user, common.ErrDataNotFound
		}
	}
	// get user by id.
	user.Id = userId
	// call user.GetUserById() method.
	// if error is not nil, return common.ErrDataNotFound error.
	if err := s.dbstore.User.GetUserById(ctx, &user); err != nil {
		// check if err contains.
		if strings.Contains(err.Error(), "no rows") {
			return user, common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return user, err
	}
	// return user.
	return user, nil
}

func (s *UserService) UpdateRole(ctx context.Context, userId int64, roleName string) error {
	// Check if current user has admin role. or current user is one of super admin users.
	currentUser := ctx.Value(common.CurrentUserKey).(types.User)
	if currentUser.Role.Name != "admin" {
		return common.ErrIncorrectDataReceived
	}
	// Get user by passed id.
	var userToBeUpdated types.User
	userToBeUpdated.Id = userId
	if err := s.dbstore.User.GetUserById(ctx, &userToBeUpdated); err != nil {
		return common.ErrIncorrectDataReceived
	}
	// Check if user to be updated in super admin user, if so return an Incorrect data error.
	if slices.Contains(utils.GetSuperUsers(), userToBeUpdated.Email) {
		return common.ErrIncorrectDataReceived
	}
	// Update user role to passed roleName.

	if err := s.dbstore.User.UpdateUserRole(ctx, &userToBeUpdated, roleName); err != nil {
		// TODO: check the error type if it containes role deosnt exist.
		// Then return ErrIncorrectDataReceived error.
		if strings.Contains(err.Error(), "no rows") {
			return common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return err

	}
	// Return nil if update is successful.
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, email string) (string, error) {
	var user types.User
	user.Email = email
	if err := s.dbstore.User.GetUserByEmail(ctx, &user); err != nil {
		return "", common.ErrIncorrectDataReceived
	}
	token, err := s.jwtService.GeneratePasswordToken(user.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		// TODO: To be updated to internal server error.
		return "", common.ErrDataNotFound
	}
	// TODO: Send mail with token to update password.
	return token, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, token, newPass string) error {
	// create a empity user variable of type types.User.
	var user types.User
	// Parse token to get userId.
	userId, err := s.jwtService.ParsePasswordToken(token)
	if err != nil {
		// check the error message and return error accordingly.
		if strings.Contains(err.Error(), "token is expired") {
			return common.ErrIncorrectDataReceived
		} else if strings.Contains(err.Error(), "signature is invalid") {
			return common.ErrDataNotFound
		}
	}
	user.Id = userId
	// Assign password to user.Password.Text
	user.Password.Text = newPass
	// Genrate Hash of new password.
	user.Password.HashPassword()
	// Store new password in database.
	if err := s.dbstore.User.UpdateUserPassword(ctx, &user); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return err

	}
	// Return nil if update is successful.
	return nil
}

func (s *UserService) Verify(ctx context.Context, token string) error {
	var user types.User
	// Parse the token. to get the user id.
	// Need to create ParseseVerifyToken function in jwtService to return user id.
	userId, err := s.jwtService.ParseVerifyToken(token)
	if err != nil {
		// check the error message and return error accordingly.
		if strings.Contains(err.Error(), "token is expired") {
			return common.ErrIncorrectDataReceived
		} else if strings.Contains(err.Error(), "signature is invalid") {
			return common.ErrDataNotFound
		}
		return err
	}
	user.Id = userId
	if err := s.dbstore.User.GetUserById(ctx, &user); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return common.ErrDataNotFound
		}
		log.Printf("err: %s", err.Error())
		return err
	}
	// Check if the user is verified or not.
	if user.IsVerified {
		// returns from function is user is already verified.
		return nil
	}
	// If user is not verified then call the Verify function from dbstore.
	if err := s.dbstore.User.Verify(ctx, &user); err != nil {
		// Not sure what kind of error will be returned by database.
		return err
	}
	return nil
}

func (s *UserService) SendVerification(ctx context.Context) (string, error) {
	// Check if current user has admin role. or current user is one of super admin users.
	currentUser := ctx.Value(common.CurrentUserKey).(types.User)
	token, err := s.jwtService.GenerateVerifyToken(currentUser.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		// TODO: To be updated to internal server error.
		return "", common.ErrDataNotFound
	}
	// TODO: Send mail with token to update password.
	return token, nil
}

func (s *UserService) AddEditor(ctx context.Context, token string) error {
	// NOTE: Token will be passed on by the caller,
	// NOTE: Token will contain the information on the creater's id.
	// NOTE: CurrentUser will be treated as editor.
	// Get current user from ctx.
	user := ctx.Value(common.CurrentUserKey).(types.User)
	createrId, err := s.jwtService.ParseEditorReqToken(token)
	if err != nil {
		// check the error message and return error accordingly.
		if strings.Contains(err.Error(), "token is expired") {
			return common.ErrIncorrectDataReceived
		} else if strings.Contains(err.Error(), "signature is invalid") {
			return common.ErrDataNotFound
		}
	}
	if err := s.dbstore.User.AddEditorById(ctx, user.Id, createrId); err != nil {
		return common.ErrIncorrectDataReceived
	}
	// Check if user id passed is same as current user id, if same return an error.
	return nil
}

func (s *UserService) SendEditorRequest(ctx context.Context, userId int64) (types.SendEditorRequestKM, error) {
	currentUser := ctx.Value(common.CurrentUserKey).(types.User)
	// Check if userId is same as current user id, if same return an error.
	if userId == currentUser.Id {
		return types.SendEditorRequestKM{}, common.ErrIncorrectDataReceived
	}
	// Check if userId exist in database. if not then return an error.
	var editorInfo types.User
	editorInfo.Id = userId
	if err := s.dbstore.User.GetUserById(ctx, &editorInfo); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return types.SendEditorRequestKM{}, common.ErrDataNotFound
		}
		// Write a better log message.
		log.Printf("SendEditorRequest error in getting userid from database.: %s", err.Error())
		return types.SendEditorRequestKM{}, err
	}
	// Add function in storage layer to check if user is already added as editor.
	var editor types.Editor
	editor.Id = userId
	if err := s.dbstore.User.GetEditorById(ctx, currentUser.Id, &editor); err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			return types.SendEditorRequestKM{}, common.ErrIncorrectDataReceived
		}
	}
	// Need to create token.
	// NOTE: Maybe here send the userid of current user.
	// NOTE: Because when editor will click on accept request
	// NOTE: The editor will add creator.
	token, err := s.jwtService.GenerateEditorReqToken(currentUser.Id)
	if err != nil {
		log.Printf("err: %s", err.Error())
		// TODO: To be updated to internal server error.
		return types.SendEditorRequestKM{}, common.ErrDataNotFound
	}
	res := types.SendEditorRequestKM{
		EditorEmail: editorInfo.Email,
		CreaterName: currentUser.Name,
		Token:       token,
	}
	return res, nil
}
