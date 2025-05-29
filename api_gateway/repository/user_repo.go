package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// "crowdfund/model"
	"github.com/rayhanadri/crowdfunding/user-service/model"
	"github.com/rayhanadri/crowdfunding/user-service/pb"
)

type UserRepository interface {
	GetUserByID(id int) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	LoginUser(user *model.User) (*model.User, error)
}

type userRepository struct {
	address string
}

func NewUserRepository(address string) UserRepository {
	return &userRepository{address: address}
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request
	req := &pb.UserIdRequest{Id: int32(id)} // Use the provided id parameter
	// Call the GetUserByID method
	res, err := client.GetUserByID(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetUserByID: %v", err)
		return nil, err
	}

	var user model.User
	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	user.ID = int(res.Id)
	user.Name = res.Name
	user.Email = res.Email
	user.Password = res.Password
	user.CreatedAt = GetCreatedAtTime
	user.UpdatedAt = GetUpdatedAtTime
	if user.ID == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	//validate user data
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(user.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request
	req := &pb.UserRequest{Id: 0, Name: user.Name, Email: user.Email, Password: user.Password}
	// Call the CreateUser method
	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CreateUser: %v", err)
		return nil, err
	}

	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	user.ID = int(res.Id)
	user.Name = res.Name
	user.Email = res.Email
	user.Password = res.Password
	user.CreatedAt = GetCreatedAtTime
	user.UpdatedAt = GetUpdatedAtTime

	return user, nil
}

func (r *userRepository) UpdateUser(user *model.User) (*model.User, error) {
	//validate user data
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(user.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request
	req := &pb.UserRequest{Id: 0, Name: user.Name, Email: user.Email, Password: user.Password}
	// Call the CreateUser method
	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("Error calling UpdateUser: %v", err)
		return nil, err
	}

	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	user.ID = int(res.Id)
	user.Name = res.Name
	user.Email = res.Email
	user.Password = res.Password
	user.CreatedAt = GetCreatedAtTime
	user.UpdatedAt = GetUpdatedAtTime

	return user, nil
}

func (r *userRepository) LoginUser(user *model.User) (*model.User, error) {
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request
	req := &pb.UserLoginRequest{Email: user.Email, Password: user.Password}
	// Call the CreateUser method
	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("Error calling UpdateUser: %v", err)
		return nil, err
	}

	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	user.ID = int(res.Id)
	user.Name = res.Name
	user.Email = res.Email
	user.Password = res.Password
	user.CreatedAt = GetCreatedAtTime
	user.UpdatedAt = GetUpdatedAtTime

	return user, nil

}
