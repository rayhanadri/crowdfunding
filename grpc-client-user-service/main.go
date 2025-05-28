package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-grpc-client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createUser(name string, email string, password string) (result string, error error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a request
	req := &pb.UserRequest{Id: 0, Name: name, Email: email, Password: password}
	// Call the CreateUser method
	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CreateUser: %v", err)
		return "", err
	}

	return res.GetMessage(), nil
}

func main() {
	// function create user
	result, err := createUser("Jhon Doe3", "jhon3@example.com", "password123")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print the response
	fmt.Printf("Response from server: %s\n", result)
}
