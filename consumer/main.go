package main

import (
	"consumer/controller"
	"consumer/driver"
	"consumer/logger"
	"consumer/utils"
	"net/http"
	"os"
)

func init() {

	//Initialize the log configurations
	logger.LoggerConfiguration()

	//Initialize the application configurations
	applicationPath, _ := os.Getwd()
	utils.Configuration(applicationPath)

	// Initialize the database connection instance
	driver.NewPostgres()

	// Initialize the consumer
	go controller.ConsumeData()
}

func main() {
	// Initialize PostgreSQL connection (if required for other routes)
	postgres := driver.GetDatabaseConnection()

	// Set up routes
	router := controller.SetupRoutes(postgres)

	// Start the server
	logger.Log.Infoln("Consumer app is running on http://localhost:8081")
	logger.Log.Fatalln(http.ListenAndServe(":8081", router))
}
