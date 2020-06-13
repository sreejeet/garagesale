package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sreejeet/garagesale/internal/platform/conf"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/schema"
)

func main() {

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("Error generating usage for database: %s\n", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("Error parsing database configuration: %s\n", err)
	}

	// Open database connection for performing actions
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatalf("error connecting to database: %s\n", err)
	}
	defer db.Close()

	switch cfg.Args.Num(0) {

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
