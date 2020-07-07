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

	_ "net/http/pprof" // Register the pprof handlers

	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/platform/conf"
	"github.com/sreejeet/garagesale/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {

	// Created the logger object
	log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			Debug           string        `conf:"default:localhost:6000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User     string `conf:"default:postgres"`
			Password string `conf:"default:postgres,noprint"`
			Host     string `conf:"default:localhost"`
			Name     string `conf:"default:postgres"`
			// Always enable TLS on live systems
			// Currently set to true for convenience
			DisableTLS bool `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing configuration")
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
		return errors.Wrap(err, "opening database")
	}
	defer db.Close()

	// Start Debug Service
	// Route '/debug/pprof' was added to the default mux by importing the net/http/pprof package.
	// Not concerned with shutting this down when the application is shutdown.
	go func() {
		log.Println("debug service listening on", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Println("debug service closed", err)
	}()

	// Create api as a http.Server
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(db, log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// A channel to listen for errors from the server.
	// A buffer is used so that the goroutine can safely exit
	// if we fail to collect the error.
	serverErrors := make(chan error, 1)

	// Here we start the server for the (micro)service
	go func() {
		log.Printf("Server started at %s\n", api.Addr)
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
		return errors.Wrap(err, "serving")

	case <-shutdown:
		log.Printf("Shutting down service")

		// Deadline for finishing any outstanding requests
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("Could not gracefully shut down server in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}
		if err != nil {
			return errors.Wrap(err, "failed stopping server gracefully")
		}
	}

	return nil
}
