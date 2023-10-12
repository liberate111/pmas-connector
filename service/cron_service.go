package service

import (
	"app-connector/config"
	"app-connector/controller"
	"app-connector/logger"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func Cronjob() {
	localTime, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		logger.Logger.Error("crontask load local time", "error", err.Error())
		os.Exit(1)
	}
	c := cron.New(cron.WithLocation(localTime))
	c.Start()

	site := config.Config.Site

	for _, v := range site {
		v := v
		log.Println("site v", v.Name)
		c.AddFunc(config.Config.App.Schedule, func() {
			log.Println("site v add func", v.Name)
			// UpdateBySite(v)
		})
	}

	done := make(chan bool, 1)
	controller.GracefulShutdown(done)
	<-done
	defer c.Stop()
}
