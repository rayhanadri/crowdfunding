package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/rayhanadri/crowdfunding/user-service/config"
	"github.com/rayhanadri/crowdfunding/user-service/model"
	"github.com/rayhanadri/crowdfunding/user-service/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) GetUserByID(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	// Extract the ID from the request
	id := req.GetId()

	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Create a user response
	response := &pb.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *UserService) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := &model.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	//validate user data
	if user.Name == "" || user.Email == "" || user.Password == "" {
		err := errors.New("name, email, and password are required")
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}
	if len(user.Password) < 6 {
		err := errors.New("password must be at least 6 characters long")
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}
		return response, err
	}

	userPass := user.Password
	userPassHash, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}
	user.Password = string(userPassHash)

	if err := config.DB.Omit("id").Create(user).Error; err != nil {
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Last(user).Error; err != nil {
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}

	// Create a user response
	response := &pb.UserResponse{
		Message:   "User created successfully",
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *UserService) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := &model.User{
		ID:       int(req.GetId()),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	userPass := user.Password
	userPassHash, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(userPassHash)

	if err := config.DB.Model(user).Updates(user).Error; err != nil {
		return nil, err
	}

	// Create a user response
	response := &pb.UserResponse{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *UserService) LoginUser(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// get pass from database and compare with user input
	var userDb model.User
	if err := config.DB.Where("email = ?", req.GetEmail()).First(&userDb).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Create a user response
	response := &pb.UserResponse{
		Id:        int32(userDb.ID),
		Name:      userDb.Name,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userDb.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}
