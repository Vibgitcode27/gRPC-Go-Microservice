package main

import (
	"flag"
	"fmt"
	"grpc/psm"
	"grpc/service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	fmt.Println("Starting on port", *port)

	laptopServer := service.NewLaptopService(service.NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()

	psm.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}

}
