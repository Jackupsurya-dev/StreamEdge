package controller

import (
	"consumer/driver"
	"consumer/logger"
	"consumer/rabbitmq"
	"consumer/redis"
	"consumer/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/spf13/viper"
)

var clients = make(map[chan driver.User]struct{})
var mu sync.Mutex

func GetUsersHandler(db *driver.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse filters from the query parameters
		filters := map[string][]string{}
		for key, values := range r.URL.Query() {
			filters[key] = values
		}

		// Fetch filtered users from the database
		users, err := db.GetFilteredUsers(filters)
		if err != nil {
			logger.Log.Errorln("Failed to fetch users:", err)
			http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the users in JSON format
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logger.Log.Errorln("Failed to encode users to JSON:", err)
			http.Error(w, "Failed to encode users to JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func StreamUsers(db *driver.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		clientChan := make(chan driver.User, 10)
		mu.Lock()
		clients[clientChan] = struct{}{}
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(clients, clientChan)
			mu.Unlock()
			close(clientChan)
		}()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		for user := range clientChan {
			jsonData, _ := json.Marshal(user)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	}
}

func ConsumeData() {
	consumer, err := rabbitmq.NewConsumer(viper.GetString("rabbitmq_url"), "csv_queue")
	if err != nil {
		logger.Log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}
	defer consumer.Close()

	dbConn := driver.GetDatabaseConnection()
	redisClient := redis.NewRedis(viper.GetString("redis_address"))

	consumer.Consume("csv_queue", func(msg string) {
		// Decrypt the message
		encryptionKey := viper.GetString("encryption_key")
		decryptedData, err := utils.Decrypt(msg, encryptionKey)
		if err != nil {
			logger.Log.Errorf("Failed to decrypt message: %v\n", err)
			return
		}
		// Parse the message into the User struct
		user := driver.User{}
		if err := json.Unmarshal([]byte(decryptedData), &user); err != nil {
			logger.Log.Errorf("Failed to parse message: %v\n", err)
			return
		}

		// Send user data to SSE clients
		broadcastUser(user) 

		// Insert the user data into the database
		err = dbConn.InsertUser(&user)
		if err != nil {
			logger.Log.Errorf("Failed to insert user into database: %v\n", err)
			return
		}

		// Save the user data into Redis
		err = redisClient.SaveUser(user.ID, user)
		if err != nil {
			logger.Log.Errorf("Failed to save user in Redis: %v\n", err)
			return
		}

		logger.Log.Infof("Successfully inserted user: %+v\n", user)

	})
}

// Broadcast new user data to all connected SSE clients
func broadcastUser(user driver.User) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		select {
		case client <- user:
		default:
			close(client)
			delete(clients, client)
		}
	}
}
