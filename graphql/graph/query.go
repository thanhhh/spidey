package graph

import (
	"context"
	"log"
)

func (s *GraphQLServer) Query_accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]Account, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOutInSecond)
	defer cancel()

	result := []Account{}

	if id != nil && len(*id) > 0 {
		account, err := s.accountClient.GetAccount(ctx, *id)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		result = append(result, Account{
			ID:   account.ID,
			Name: account.Name,
		})

		return result, nil
	}

	skip, take := uint64(0), uint64(0)

	if pagination != nil {
		skip, take = uint64(*pagination.Skip), uint64(*pagination.Take)
	}

	accounts, err := s.accountClient.GetAccounts(ctx, skip, take)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, a := range accounts {
		result = append(result, Account{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return result, err
}

func (s *GraphQLServer) Query_products(ctx context.Context, pagination *PaginationInput, query *string, ids []string) ([]Product, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOutInSecond)
	defer cancel()

	queryValue := ""

	skip, take := uint64(0), uint64(0)

	if pagination != nil {
		skip, take = uint64(*pagination.Skip), uint64(*pagination.Take)
	}

	if query != nil {
		queryValue = *query
	}

	products, err := s.catalogClient.GetProducts(ctx, skip, take, ids, queryValue)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := []Product{}

	for _, p := range products {
		result = append(result,
			Product{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
	}

	return result, nil
}
