package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rayhanadri/crowdfunding/user-service/config"
	"github.com/rayhanadri/crowdfunding/user-service/pb"
	"github.com/rayhanadri/crowdfunding/user-service/service"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()
	// Connect to the database
	config.Connect()

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the UserService with the gRPC server
	pb.RegisterUserServiceServer(grpcServer, &service.UserService{})

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Server is running on port :50051...")

	// Serve gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
