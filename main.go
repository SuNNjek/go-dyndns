package main

import (
	"context"
	"errors"
	"log"
)

func main() {
	ctx, disposeCtx := wrapInterruptContext(context.Background())
	defer disposeCtx()

	client, err := Init()
	if err != nil {
		log.Fatalln(err)
	}

	// Run the client until we get an interrupt or an error
	if err := client.Run(ctx); errors.Is(err, context.Canceled) {
		// No need to return an error code in case of requested cancellation
		log.Println("Shutdown requested, shutting down...")
	} else if err != nil {
		log.Fatalln(err)
	}
}
