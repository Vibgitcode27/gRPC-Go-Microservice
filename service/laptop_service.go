package service

import (
	"context"
	"errors"
	"fmt"
	"grpc/psm"
	"log"

	// "grpc/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopService struct {
	psm.UnimplementedLaptopServiceServer
	Store LaptopStore
}

func NewLaptopService(store LaptopStore) *LaptopService {
	return &LaptopService{
		Store: store,
	}
}

func (server *LaptopService) CreateLaptop(ctx context.Context, req *psm.CreateLaptopRequest) (*psm.CreateLaptopResponse, error) {

	laptop := req.GetLaptop()
	fmt.Printf("Receive a create-laptop request with id: %s\n", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	err := server.Store.Save(laptop)

	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			// return nil, status.Errorf(code, "Laptop with ID %s already exists", laptop.Id)
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save %s beacuse of %v", laptop.Id, err)
	}

	fmt.Println("Laptop saved successfully")
	return &psm.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}

func (server *LaptopService) SearchLaptop(req *psm.SearchLaptopRequest, stream psm.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("Receive a search-laptop request with filter: %v\n", filter)

	err := server.Store.Search(filter, func(laptop *psm.Laptop) error {
		res := &psm.SearchLaptopResponse{
			Laptop: laptop,
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		log.Printf("Sent all laptop Id with filter: %v\n", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}
