package main

import (
	"flag"
	"log"
	"os"

	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/schema"
)

func main() {
	flag.Parse()

	// Open database connection for performing actions
	db, err := database.Open()
	if err != nil {
		log.Fatalf("error connecting to database: %s\n", err)
	}
	defer db.Close()

	switch flag.Arg(0) {

	case "migrate":
		// Migrate database schema
		if err := schema.Migrate(db); err != nil {
			log.Println("Error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Completed migration")
		return

	case "seed":
		// Seeding database
		if err := schema.Seed(db); err != nil {
			log.Println("Error seeding database", err)
			os.Exit(1)
		}
		log.Println("Completed seeding database")
		return
	}

}
