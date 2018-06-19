package order

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/thanhhh/spidey/account"
	"github.com/thanhhh/spidey/catalog"
	"github.com/thanhhh/spidey/order/pb"

	"google.golang.org/grpc"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenRGPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)

	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)

	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &grpcServer{
		s, accountClient, catalogClient})

	return server.Serve(listener)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	productIds := []string{}

	for _, p := range r.Products {
		productIds = append(productIds, p.ProductId)
	}

	orderProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIds, "")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	products := []OrderedProduct{}

	for _, p := range orderProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}

		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}

	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products,
			&pb.Order_OrderProduct{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}
func (s *grpcServer) GetOrdersForAccount(
	ctx context.Context,
	r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {

	orders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	productFoundMap := make(map[string]*catalog.Product)
	for _, o := range orders {
		for _, p := range o.Products {
			if _, exist := productFoundMap[p.ID]; !exist {
				productFoundMap[p.ID] = nil
			}
		}
	}

	productIDs := []string{}

	for id := range productFoundMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, p := range products {
		productFoundMap[p.ID] = &p
	}

	response := &pb.GetOrdersForAccountResponse{}
	response.Orders = []*pb.Order{}

	for _, o := range orders {
		orderResp := &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
		}
		orderResp.CreatedAt, _ = o.CreatedAt.MarshalBinary()
		orderResp.Products = []*pb.Order_OrderProduct{}

		for _, op := range o.Products {
			p := productFoundMap[op.ID]

			opResp := &pb.Order_OrderProduct{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    op.Quantity,
			}

			orderResp.Products = append(orderResp.Products, opResp)
		}

		response.Orders = append(response.Orders, orderResp)
	}

	return response, nil
}
