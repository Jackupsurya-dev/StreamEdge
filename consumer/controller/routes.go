package controller

import (
	"consumer/driver"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *driver.Postgres) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsersHandler(db)).Methods("GET")
	router.HandleFunc("/users/stream", StreamUsers(db)).Methods("GET")
	return router
}
