package main

import (
	"github.com/JneiraS/GotoServ/internal/api"
	"github.com/JneiraS/GotoServ/pkg/utils"
)

func main() {

	utils.LoadEnvironmentVariables()
	utils.CreatJsonFromCsv()
	router := api.NewRouter("assignement_fcb")
	api.StartServer(router)
}
