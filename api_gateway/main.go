package main

import (
	"crowdfund/config" // Import the config package
	"crowdfund/route"  // Import the route package
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
