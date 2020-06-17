package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/internal/platform/conf"
	"github.com/sreejeet/garagesale/internal/platform/database"
	"github.com/sreejeet/garagesale/internal/schema"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}

func run() error {

	var cfg struct {
		DB struct {
			User     string `conf:"default:postgres"`
			Password string `conf:"default:postgres,noprint"`
			Host     string `conf:"default:localhost"`
			Name     string `conf:"default:postgres"`
			// Always enable TLS on live systems
			// Currently set to true for convenience
			DisableTLS bool `conf:"default:true"`
		}
		Args conf.Args
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
		return errors.Wrap(err, "parsing config")
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
		return errors.Wrap(err, "connecting to database")
	}
	defer db.Close()

	switch cfg.Args.Num(0) {

	case "migrate":
		// Migrate database schema
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		log.Println("Completed migration")
		return nil

	case "seed":
		// Seeding database
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		log.Println("Completed seeding database")
		return nil
	}

	return nil
}
