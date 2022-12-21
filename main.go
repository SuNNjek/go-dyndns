package main

import (
	"context"
	"errors"
	"go-dyndns/client"
	"go-dyndns/log"
	golog "log"
	"os"
	"os/signal"
	"syscall"
)

// Listen to SIGINT (Ctrl+C) and SIGTERM (docker stop) signals
var cancelSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), cancelSignals...)
	defer cancel()

	logger, err := log.CreateLogger()
	if err != nil {
		golog.Fatalln(err)
	}

	logger.Info("Started go-dyndns")

	cl, err := client.CreateClient(logger)
	if err != nil {
		logger.Fatal("Failed to create client: %v", err)
	}

	// Run the client until we get an interrupt or an error
	if err := cl.Run(ctx); errors.Is(err, context.Canceled) {
		// No need to return an error code in case of requested cancellation
		logger.Info("Shutdown requested, shutting down...")
	} else if err != nil {
		logger.Fatal("Unhandled error: %v", err)
	}
}
