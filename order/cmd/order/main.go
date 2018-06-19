package main

import (
	"log"
	"time"

	"github.com/tinrab/retry"

	"github.com/thanhhh/spidey/order"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})

	defer r.Close()

	s := order.NewService(r)

	log.Fatal(order.ListenRGPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
