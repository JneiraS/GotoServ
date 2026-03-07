package main

import (
	"log"

	"github.com/JneiraS/GotoServ/internal/api"
	"github.com/JneiraS/GotoServ/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentVariables()
	utils.CreatJsonFromCsv()
	router := api.NewRouter("assignement_fcb.json")
	api.StartServer(router)
}

func loadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}
}
