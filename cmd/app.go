package main

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/config"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/security"
	"AvitoTechPVZ/service"
	"fmt"

	"log"
	"net/http"
)

func main() {
	// Get Cfg
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error loading App config: %v", err)
	}

	// Open DB conn
	db, err := repo.NewDbConn(cfg.Db)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer repo.CloseDbConn(db)

	// Create service
	repo := repo.NewRepo(db)
	sc := security.NewSecurityController(cfg.App.Secret)
	s := service.NewDefaultAPIServicerImpl(repo, sc)
	controller := dto.NewDefaultAPIController(s)

	// Register routes
	router := s.NewRouter(controller)

	// Starting server
	// TODO: gracefull shutdown
	addr := fmt.Sprintf(":%v", cfg.App.Port)
	log.Fatal(http.ListenAndServe(addr, router)) // Fatal убивает процесс и все деферы идут лесом
}
