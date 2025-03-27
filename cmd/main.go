package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heshanu/go-service/pkg/config"
	"github.com/heshanu/go-service/pkg/routes"
)

func main() {
	config.Connect()
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)

	http.Handle("/", r)

	// Load environment variables from .env file
	http.ListenAndServe(":8081", r)
	fmt.Printf(":Application is running on:8081")

}
