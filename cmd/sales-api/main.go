package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"

	_ "expvar"         // Register expvars handlers
	_ "net/http/pprof" // Register pprof handlers

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/cmd/sales-api/internal/handlers"
	"github.com/sreejeet/garagesale/internal/platform/auth"
	"github.com/sreejeet/garagesale/internal/platform/conf"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"go.opencensus.io/trace"
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
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api/v2/spans"`
			Service     string  `conf:"default:sales-api"`
			Probability float64 `conf:"default:1"`
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
	log.Print("Started service")
	defer log.Print("Ended service")

	// Initialize authentication support
	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

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

	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()

	// Start Debug Service
	//
	// Route '/debug/pprof' was added to the default mux by importing the net/http/pprof package.
	// Route '/debug/vars' was added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.
	go func() {
		log.Println("debug service listening on", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Println("debug service closed", err)
	}()

	// Another channel to receive OS signals like SIGINT or SIGTERM.
	// The signal package requires this channel to be buffered.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Create api as a http.Server
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(shutdown, db, log, authenticator),
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

	// Using the switch case construct, or in case of Go,
	// the select construct to block main func till shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "serving")

	case <-shutdown:
		log.Printf("Shutting down service")

	case sig := <-shutdown:
		log.Printf("main : received %v : starting shutdown", sig)

		// Deadline for finishing any outstanding requests
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("Could not gracefully shut down server in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}
		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("critical or unexpected error issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "failed stopping server gracefully")
		}
	}

	return nil
}

// createAuth ceates an authenticator. It reads a private key file, looks up the the
// public key for it, then creats an authenticator using the algorithm provided for signing.
func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

// registerTracer is used to register a tracer for a particular service
func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {

	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
