package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/platform/conf"
	"github.com/sreejeet/garagesale/internal/platform/database"
)

func main() {

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	const serveURL string = "localhost:8000"

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	// Basic logging
	log.Printf("Started service")
	defer log.Print("Ended service")

	// Start database
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
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
