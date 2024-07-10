package main

import (
	"context"
	"flag"
	"fmt"
	"grpc/psm"
	"grpc/sample"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "The server address in the format of host:port")
	flag.Parse()
	fmt.Println("Dial server on address", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	laptopClient := psm.NewLaptopServiceClient(conn)

	fmt.Println("laptopClient", laptopClient)

	laptop := sample.Laptop()

	laptop.Id = "invalid"

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
