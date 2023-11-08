package controller

import (
	"app-connector/logger"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(done chan bool) {
	// Gracefully Shutdown
	// Make channel listen for signals from OS
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-c
		logger.Info("event: shutdown_app, status: process, msg: Application shutdown...")
		CloseDB()
		done <- true
	}()
}
