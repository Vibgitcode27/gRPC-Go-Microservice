package service

import (
	"context"
	"errors"
	"fmt"
	"grpc/psm"

	// "grpc/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopService struct {
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
