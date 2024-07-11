package main

import (
	"context"
	"flag"
	"fmt"
	"grpc/psm"
	"grpc/sample"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateRandomLaptop(laptopClient psm.LaptopServiceClient) {
	fmt.Println("laptopClient", laptopClient)

	laptop := sample.Laptop()

	// laptop.Id = "invalid"

	req := &psm.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			fmt.Println("Laptop already exists")
		} else {
			log.Fatal("Cannot create laptop", err)
		}
		return
	}

	log.Printf("Created laptop with id: %s", res.Id)
}

func main() {
	serverAddress := flag.String("address", "", "The server address in the format of host:port")
	flag.Parse()
	fmt.Println("Dial server on address", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	laptopClient := psm.NewLaptopServiceClient(conn)

	for i := 0; i < 5; i++ {
		CreateRandomLaptop(laptopClient)
	}

	filter := &psm.Filter{
		MaxPriceInr: 150000,
		MinCpuCores: 2,
		MinCpuGhz:   2.0,
		Ram:         &psm.Memory{Value: 8, Unit: psm.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
}

func searchLaptop(laptopClient psm.LaptopServiceClient, filter *psm.Filter) {
	fmt.Println("Search laptop with filter", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &psm.SearchLaptopRequest{
		Filter: filter,
	}

	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("Cannot search laptop", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("")
			fmt.Printf("")
			log.Printf("End of the file")
			return
		}

		if err != nil {
			log.Fatal("Cannot receive response", err)
		}

		laptop := res.GetLaptop()
		log.Printf("Found laptop with id: %s", laptop.GetId())
		log.Printf("Laptop brand: %s", laptop.GetBrand())
		log.Printf("Laptop name: %s", laptop.GetName())
		log.Printf("Laptop CPU cores: %d", laptop.GetCpu().GetCores())
		log.Printf("Laptop CPU ghz: %f", laptop.GetCpu().GetMinGhz())
		log.Printf("Laptop RAM: %d %s", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
	}
}
