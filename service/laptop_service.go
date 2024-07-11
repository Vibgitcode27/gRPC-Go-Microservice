package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"grpc/psm"
	"io"
	"log"

	// "grpc/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxImageSize = 1 << 20

type LaptopService struct {
	psm.UnimplementedLaptopServiceServer
	laptopStore LaptopStore
	imageStore  ImageStore
}

func NewLaptopService(laptopStore LaptopStore, imageStore ImageStore) *LaptopService {
	return &LaptopService{laptopStore: laptopStore, imageStore: imageStore}
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

	err := server.laptopStore.Save(laptop)

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

	err := server.laptopStore.Search(filter, func(laptop *psm.Laptop) error {
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

func (server *LaptopService) UploadImage(stream psm.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info: %v", err)
	}

	laptopID := req.GetImage().GetImageId()
	imageType := req.GetImage().GetImageType()

	log.Printf("Receive an upload-image request for laptop %s with image type %s\n", laptopID, imageType)

	laptop, err := server.laptopStore.Find(laptopID)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot find laptop: %v", err)
	}

	if laptop == nil {
		return status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist", laptopID)
	}

	imageDate := bytes.Buffer{}
	imageSize := 0

	for {
		log.Print("Waiting to receive more data")

		req, err := stream.Recv()

		if err == io.EOF {
			log.Print("No more data")
			break
		}

		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize >= MaxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, MaxImageSize)
		}

		_, err = imageDate.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}

	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageDate)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to store: %v", err)
	}

	res := &psm.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)

	if err != nil {
		return status.Errorf(codes.Internal, "cannot send response: %v", err)
	}

	log.Printf("Image with ID %s and size %d has been saved successfully", imageID, imageSize)

	return nil
}
