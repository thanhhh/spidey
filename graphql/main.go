package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/thanhhh/spidey/graphql/graph"
	"github.com/vektah/gqlgen/handler"
)

type Config struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	s, err := graph.NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(s)))
	http.Handle("/playground", handler.Playground("Spidey", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
