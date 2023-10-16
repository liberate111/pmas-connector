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
		logger.Error("connect to db", "error", err.Error())
		os.Exit(1)
	}
}

func main() {
	if config.Config.App.InitialDB {
		service.InitTable()
	}
	service.Cronjob()
}
