package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// wrapInterruptContext wraps the given context, so that it gets cancelled on SIGINT or SIGTERM
func wrapInterruptContext(ctx context.Context) (context.Context, func()) {
	wrappedCtx, cancel := context.WithCancel(ctx)

	// Listen to SIGINT (Ctrl+C) and SIGTERM (docker stop) signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, syscall.SIGTERM)

	// Function that stops listening for signals. Needs to be called by the caller
	disposeFunc := func() {
		signal.Stop(signalChan)
		cancel()
	}

	go func() {
		select {
		// When signal is received, cancel the context
		case <-signalChan:
			cancel()
		case <-wrappedCtx.Done():
		}
	}()

	// Return the context and the disposeFunc, so it can be called by the calling function to stop listening for the signal
	return wrappedCtx, disposeFunc
}
