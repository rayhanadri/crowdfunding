package main

import (
	"context"
	"fmt"
	"log"
	"time"

	donation_pb "github.com/rayhanadri/crowdfunding/donation-service/pb"
	"github.com/rayhanadri/crowdfunding/user-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func GetUserByID(id int32) (result string, error error) {
	// conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	address := "user-service-273575294549.asia-southeast2.run.app:443"
	// creds := credentials.NewClientTLSFromCert(nil, "") // use system root CAs
	// conn, err := grpc.Dial("https://user-service-273575294549.asia-southeast2.run.app:443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("Did not connect: %v", err)
	// }
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return "", err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a request
	req := &pb.UserIdRequest{Id: id} // Use the provided id parameter
	// Call the GetUserByID method
	res, err := client.GetUserByID(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetUserByID: %v", err)
		return "", err
	}

	result = fmt.Sprintf("ID: %d, Name: %s, Email: %s, CreatedAt: %s, UpdatedAt: %s",
		res.GetId(), res.GetName(), res.GetEmail(), res.GetCreatedAt(), res.GetUpdatedAt())

	return result, nil
}

func LoginUser(email string, password string) (result string, error error) {
	// conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	address := "user-service-273575294549.asia-southeast2.run.app:443"
	// creds := credentials.NewClientTLSFromCert(nil, "") // use system root CAs
	// conn, err := grpc.Dial("https://user-service-273575294549.asia-southeast2.run.app:443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("Did not connect: %v", err)
	// }
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return "", err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a request
	req := &pb.UserLoginRequest{Email: email, Password: password} // Use the provided id parameter
	// Call the LoginUser method
	res, err := client.LoginUser(ctx, req)
	if err != nil {
		log.Fatalf("Error calling LoginUser: %v", err)
		return "", err
	}

	result = fmt.Sprintf("ID: %d, Name: %s, Email: %s, CreatedAt: %s, UpdatedAt: %s",
		res.GetId(), res.GetName(), res.GetEmail(), res.GetCreatedAt(), res.GetUpdatedAt())

	return result, nil
}

func GetDonationByID(id int32) (result string, error error) {
	/// conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// address := "donation-service-273575294549.asia-southeast2.run.app:443" // or "donation-service-273575294549.asia-southeast2.run.app:443"
	address := "localhost:50051" // or "donation-service-273575294549.asia-southeast2.run.app:443"
	// creds := credentials.NewClientTLSFromCert(nil, "") // use system root CAs
	// conn, err := grpc.Dial("https://user-service-273575294549.asia-southeast2.run.app:443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("Did not connect: %v", err)
	// }
	conn, err := grpc.Dial(
		address,
		// grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()), // for secure TLS
	)

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return "", err
	}

	defer conn.Close()

	// Create a new client
	client := donation_pb.NewDonationServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a request
	req := &donation_pb.DonationIdRequest{Id: int32(id)} // Use the provided id parameter
	// Call the GetUserByID method
	res, err := client.GetDonationByID(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetDonationByID: %v", err)
		return "", err
	}

	result = fmt.Sprintf("ID: %d, Name: %d, Amount: %.2f, Message: %s,CreatedAt: %s, UpdatedAt: %s",
		res.GetId(), res.CampaignId, res.GetAmount(), res.GetMessageText(), res.GetCreatedAt(), res.GetUpdatedAt())

	return result, nil
}

func main() {
	// function create user
	// result, err := createUser("Jhon Doe4", "jhon4@example.com", "password123")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
	result, err := GetUserByID(3)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response from server: %s\n", result)

	result, err = LoginUser("jhon4@example.com", "password123")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Print the response
	fmt.Printf("Response from server: %s\n", result)

	result, err = GetDonationByID(3)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response from server: %s\n", result)
}
