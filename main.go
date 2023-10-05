package main

import (
	"app-connector/config"
	"app-connector/controller"
	"app-connector/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	var err error
	config.ReadConfig()
	logger.InitLog()
	controller.InitClient()
	err = controller.ConnectDB()
	if err != nil {
		os.Exit(1)
	}

	if config.Config.App.InitialDB {
		// InitTable()
	}
}

func main() {
	time.Sleep(1 * time.Minute)
	// cronjob()
}

func gracefulShutdown(done chan bool) {
	// Gracefully Shutdown
	// Make channel listen for signals from OS
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-c
		logger.Logger.Info("Application shutdown...")
		controller.CloseDB()
		done <- true
	}()
}

func cronjob() {
	localTime, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		logger.Logger.Error("crontask load local time", "error", err.Error())
		os.Exit(1)
	}
	c := cron.New(cron.WithLocation(localTime))

	// ===================== ALL SITE ======================================
	// chn-c2
	// c.AddFunc("@midnight", chn_c2)
	// for test
	// c.AddFunc("@every 30s", updateBySite)s

	// other site
	// ===================== ALL SITE ======================================

	c.Start()
	done := make(chan bool, 1)
	gracefulShutdown(done)
	<-done
	defer c.Stop()
}
