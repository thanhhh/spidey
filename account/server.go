package account

import (
	"context"
	"fmt"
	"net"

	"github.com/thanhhh/spidey/account/pb"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{s})
	reflection.Register(server)

	return server.Serve(listen)
}

func (s *grpcServer) PostAccount(ctx context.Context, request *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {

	a, err := s.service.PostAccount(ctx, request.Name)

	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, request *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, request.Id)

	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, request *pb.GetAccountsRequest) (*pb.GetAccountsReponse, error) {
	accounts, err := s.service.GetAccounts(ctx, request.Skip, request.Take)

	if err != nil {
		return nil, err
	}

	var accountResponses []*pb.Account

	for _, a := range accounts {
		accountResponses = append(
			accountResponses,
			&pb.Account{
				Id:   a.ID,
				Name: a.Name,
			},
		)
	}

	return &pb.GetAccountsReponse{
		Accounts: accountResponses,
	}, nil
}
