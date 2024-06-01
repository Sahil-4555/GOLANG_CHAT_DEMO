package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sahil-4555/mvc/shared/log"
)

func GracefulStop(log log.ILogger, callback func(context.Context) error) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulStop

	log.Info("", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := callback(ctx); err != nil {
		log.Fatal("", "Server forced to shutdown:", err)
	}

	log.Info("", "Server exiting")
}
