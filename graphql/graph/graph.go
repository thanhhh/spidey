package graph

import (
	"errors"
	"log"
	"time"

	"github.com/thanhhh/spidey/account"
	"github.com/thanhhh/spidey/catalog"
	"github.com/thanhhh/spidey/order"
)

var (
	ErrOrderQuantityInvalid = errors.New("Order Quantity is invalid")
	TimeOutInSecond         = 3 * time.Second
)

type GraphQLServer struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountURL, catalogURL, orderURL string) (*GraphQLServer, error) {
	accountClient, err := account.NewClient(accountURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	orderClient, err := order.NewClient(orderURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &GraphQLServer{accountClient, catalogClient, orderClient}, nil
}
