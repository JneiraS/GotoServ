package main

import (
	"fmt"

	"github.com/JneiraS/GotoServ/internal/api"
	"github.com/JneiraS/GotoServ/pkg/utils"
)

func main() {
	utils.LoadEnvironmentVariables()
	utils.CreatJsonFromCsv()
	router := api.NewRouter("assignement_fcb.json")
	api.StartServer(router)
	code, _ := utils.GenerateCurrentTOTP("JHTSW6W3")
	fmt.Printf("%s", code)
}
