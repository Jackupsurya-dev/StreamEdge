package main

import (
	"net/http"
	"os"
	"producer/controller"
	"producer/logger"
	"producer/rabbitmq"
	"producer/utils"

	"github.com/spf13/viper"
)

func init() {
	//Initialize the application configurations
	applicationPath, _ := os.Getwd()
	utils.Configuration(applicationPath)

	//Initialize the log configurations
	logger.LoggerConfiguration()
}

func main() {
	// Initialize RabbitMQ producer
	producer, err := rabbitmq.NewProducer(viper.GetString("rabbitmq_url"))
	if err != nil {
		logger.Log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}
	defer producer.Close()

	// Set up routes
	router := controller.SetupRoutes(producer)

	// Start the server
	logger.Log.Infoln("producer app is running on http://localhost:8080")
	logger.Log.Fatalln(http.ListenAndServe(":8080", router))
}
