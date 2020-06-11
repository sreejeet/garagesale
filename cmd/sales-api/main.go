package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/platform/database"
)

func main() {

	const serveURL string = "localhost:8000"

	// Basic logging
	log.Printf("Started service")
	defer log.Print("Ended service")

	// Start database
	db, err := database.Open()
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}
	defer db.Close()

	productsHandler := handlers.Products{DB: db}

	// Create api as a http.Server
	api := http.Server{
		Addr:         serveURL,
		Handler:      http.HandlerFunc(productsHandler.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// A channel to listen for errors from the server.
	// A buffer is used so that the goroutine can safely exit
	// if we fail to collect the error.
	serverErrors := make(chan error, 1)

	// Here we start the server for the (micro)service
	go func() {
		log.Printf("Server started at %s\n", serveURL)
		serverErrors <- api.ListenAndServe()
	}()

	// Another channel to receive OS signals like SIGINT or SIGTERM.
	// The signal package requires this channel to be buffered.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Using the switch case construct, or in case of Go,
	// the select construct to block main func till shutdown
	select {

	case err := <-serverErrors:
		log.Fatalf("Error listening: %s\n", err)

	case <-shutdown:
		log.Printf("Shutting down service")

		// Deadline for finishing any outstanding requests
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shut down
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("Could not gracefully shut down server in %v : %v", timeout, err)
			err = api.Close()
		}
		if err != nil {
			log.Fatalf("Could not stop server gracefully : %v", err)
		}

	}
}
