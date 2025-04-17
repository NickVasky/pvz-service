package main

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/service"
	"strings"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Create service implementation
	s := &service.DefaultAPIServicerImpl{
		Users: repo.UserRepo{DB: repo.OpenDbConnection()},
	}

	// Create controller
	controller := dto.NewDefaultAPIController(s)

	// Register routes
	for _, route := range controller.Routes() {
		var h http.HandlerFunc
		if !strings.Contains(route.Pattern, "Login") {
			h = service.AuthMiddleware(route.HandlerFunc)
		} else {
			h = route.HandlerFunc
		}
		router.HandleFunc(route.Pattern, h).Methods(route.Method)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
