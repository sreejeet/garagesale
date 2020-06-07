package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	const serveURL string = "localhost:8000"

	// Basic logging
	log.Printf("Started service")
	defer log.Print("Ended service")

	// Create api as a http.Server
	api := http.Server{
		Addr:         serveURL,
		Handler:      http.HandlerFunc(ListProducts),
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

	// Another channel to recieve OS signals like SIGINT or SIGTERM.
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

// Product is a type declared for items in our garage sale
type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ListProducts is an http handler for returning
// a json list of products.
func ListProducts(w http.ResponseWriter, r *http.Request) {

	list := []Product{
		{Name: "Oil painting", Cost: 500, Quantity: 1},
		{Name: "Intel Pentium 4 CPU", Cost: 5000, Quantity: 1},
		{Name: "Fresh Pizza from 2004", Cost: 2, Quantity: 5},
	}

	// Marshalling (converting) product slice to json array
	data, err := json.Marshal(list)
	if err != nil {
		log.Print("Error parsing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		log.Print("Error writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
