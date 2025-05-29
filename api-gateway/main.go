package main

import (
	"github.com/rayhanadri/crowdfunding/api-gateway/config" // Import the config package
	"github.com/rayhanadri/crowdfunding/api-gateway/route"  // Import the route package
)

func init() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	// Initialize the database connection
	// config.Connect()
}

func main() {
	// @title Crowdfunding API
	// @version 1.0
	// @description This is a sample server for a crowdfunding API.
	// @host localhost:8080
	// @BasePath /api/v1/
	route.ExecRouter()
}
