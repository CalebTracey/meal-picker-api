package main

import (
	"food-roulette-api/internal/facade"
	"food-roulette-api/internal/routes"
	"food-roulette-api/internal/services"
	"github.com/NYTimes/gziphandler"
	config "github.com/calebtracey/config-yaml"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	configPath = "config.yaml"
)

const Port = "6080"

func main() {
	defer panicQuit()

	appConfig := config.NewFromFile(configPath)
	service, err := facade.NewService(appConfig)

	if err != nil {
		log.Panicln(err)
	}

	handler := routes.Handler{
		Service: &service,
	}

	router := handler.InitializeRoutes()

	log.Fatal(services.ListenAndServe(Port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}
