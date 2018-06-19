package graph

import (
	"context"
	"log"

	"github.com/thanhhh/spidey/order"
)

func (s *GraphQLServer) Mutation_createAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOutInSecond)
	defer cancel()

	account, err := s.accountClient.PostAccount(ctx, in.Name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}

func (s *GraphQLServer) Mutation_createProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOutInSecond)
	defer cancel()

	product, err := s.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (s *GraphQLServer) Mutation_createOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOutInSecond)
	defer cancel()

	orderProducts := []order.OrderedProduct{}

	for _, iop := range in.Products {
		if iop.Quantity <= 0 {
			return nil, ErrOrderQuantityInvalid
		}

		orderProducts = append(orderProducts,
			order.OrderedProduct{
				ID:       iop.ID,
				Quantity: uint32(iop.Quantity),
			})
	}

	order, err := s.orderClient.PostOrder(ctx, in.AccountId, orderProducts)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderResp := &Order{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		Products:   []OrderedProduct{},
	}

	for _, op := range order.Products {
		orderResp.Products = append(orderResp.Products,
			OrderedProduct{
				ID:          op.ID,
				Name:        op.Name,
				Description: op.Description,
				Price:       op.Price,
				Quantity:    int(op.Quantity),
			})
	}

	return orderResp, nil
}
