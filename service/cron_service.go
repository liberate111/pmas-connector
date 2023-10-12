package service

import (
	"app-connector/controller"
	"app-connector/logger"
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

	// schedule
	// ===================== ALL SITE ======================================
	// chn-c2
	// c.AddFunc("@midnight", service.UpdateBySite)
	// for test
	// c.AddFunc("@every 30s", updateBySite)

	// other site
	// ===================== ALL SITE ======================================

	// manual
	// reqCon := config.Config.Api.Connect
	// reqGet := config.Config.Api.GetData
	// table := config.Config.TableName
	// UpdateBySite(reqCon, reqGet, table)

	c.Start()
	done := make(chan bool, 1)
	controller.GracefulShutdown(done)
	<-done
	defer c.Stop()
}
