package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sreejeet/garagesale/schema"
)

func main() {

	const serveURL string = "localhost:8000"

	// Basic logging
	log.Printf("Started service")
	defer log.Print("Ended service")

	// Start database
	db, err := openDB()
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("Error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Completed migration")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("Error seeding database", err)
			os.Exit(1)
		}
		log.Println("Completed seeding database")
		return
	}

	service := Products{db: db}

	// Create api as a http.Server
	api := http.Server{
		Addr:         serveURL,
		Handler:      http.HandlerFunc(service.List),
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
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// Products struct holds all business logic related to products
type Products struct {
	db *sqlx.DB
}

// ListProducts is an http handler for returning
// a json list of products.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {

	list := []Product{}

	const query = `SELECT * FROM products;`

	if err := p.db.Select(&list, query); err != nil {
		log.Printf("Error retrieving product list: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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

func openDB() (*sqlx.DB, error) {
	urlQuery := url.Values{}
	urlQuery.Set("sslmode", "disable")
	urlQuery.Set("timezone", "utc")

	url := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "psotgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: urlQuery.Encode(),
	}

	return sqlx.Open("postgres", url.String())
}
