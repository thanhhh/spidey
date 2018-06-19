package main

import (
	"log"
	"time"

	"github.com/tinrab/retry"

	"github.com/kelseyhightower/envconfig"
	"github.com/thanhhh/spidey/catalog"
)

type Config struct {
	ElasticURL string `envconfig:"ELASTIC_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.ElasticURL)
		if err != nil {
			log.Println(err)
		}

		return
	})

	defer r.Close()

	log.Println("Listening on port 8080...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
