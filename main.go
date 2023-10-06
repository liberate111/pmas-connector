package main

import (
	"app-connector/config"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/service"
	"os"
)

func init() {
	config.ReadConfig()
	logger.InitLog()
	controller.InitClient()
	err := controller.ConnectDB()
	if err != nil {
		os.Exit(1)
	}

	if config.Config.App.InitialDB {
		// InitTable()
	}
}

func main() {
	service.Cronjob()
}
