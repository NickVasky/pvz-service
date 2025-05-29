package main

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/config"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/service"
	"fmt"

	"log"
	"net/http"
)

func main() {
	// Get Cfg
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create service
	db, err := repo.NewDbConn(cfg.Db)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: defer close
	repo := repo.NewRepo(db)
	s := &service.DefaultAPIServicerImpl{
		Repo: repo,
	}

	// Create controller
	controller := dto.NewDefaultAPIController(s)

	// Register routes
	//router := mux.NewRouter()
	router := service.NewRouter(controller)

	addr := fmt.Sprintf(":%v", cfg.App.Port)
	log.Fatal(http.ListenAndServe(addr, router)) // Fatal убивает процесс и все деферы идут лесом
}
