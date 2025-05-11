package main

import (
	"AvitoTechPVZ/codegen/dto"
	cfg "AvitoTechPVZ/config"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/service"
	"fmt"
	"strings"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Get Cfg
	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create service implementation
	db := repo.OpenDbConnection(config.Db)
	repo := repo.NewRepo(db)
	s := &service.DefaultAPIServicerImpl{
		Repo: repo,
	}

	// Create controller
	controller := dto.NewDefaultAPIController(s)

	// Register routes
	router := mux.NewRouter()
	for _, route := range controller.Routes() {
		var h http.HandlerFunc
		if !strings.Contains(route.Pattern, "Login") {
			h = service.AuthMiddleware(route.HandlerFunc, []string{"moderator"})
		} else {
			h = route.HandlerFunc
		}
		router.HandleFunc(route.Pattern, h).Methods(route.Method)
	}

	addr := fmt.Sprintf(":%v", config.App.Port)
	log.Fatal(http.ListenAndServe(addr, router)) // Fatal убивает процесс и все деферы идут лесом
}
