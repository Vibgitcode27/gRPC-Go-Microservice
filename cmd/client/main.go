package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"grpc/psm"
	"grpc/sample"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateRandomLaptop(laptopClient psm.LaptopServiceClient, laptop *psm.Laptop) {
	fmt.Println("laptopClient", laptopClient)

	// laptop := sample.Laptop()

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

func uploadImage(laptopClient psm.LaptopServiceClient, laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Cannot open image file", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImage(ctx)
	if err != nil {
		log.Fatal("Cannot upload image", err)
	}

	req := &psm.UploadImageRequest{
		Data: &psm.UploadImageRequest_Image{
			Image: &psm.Image{
				ImageId:   laptopID, // The image ID is the same as the laptop ID this is mistake I made in naming the field in proto file
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("Cannot send image", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Cannot read chunk to buffer", err)
		}

		req := &psm.UploadImageRequest{
			Data: &psm.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			err2 := stream.RecvMsg(nil)
			log.Fatal("Cannot send chunk to server", err, err2)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Cannot receive response from server", err)
	}

	log.Printf("Image uploaded with id: %s, size: %d", res.Id, res.Size)

}

func testUploadImage(laptopClient psm.LaptopServiceClient) {
	laptop := sample.Laptop()
	CreateRandomLaptop(laptopClient, laptop)
	uploadImage(laptopClient, laptop.GetId(), "tmp/laptop.jpg")
}

func testSearchLaptop(laptopClient psm.LaptopServiceClient) {
	for i := 0; i < 5; i++ {
		CreateRandomLaptop(laptopClient, sample.Laptop())
	}

	filter := &psm.Filter{
		MaxPriceInr: 200000,
		MinCpuCores: 2,
		MinCpuGhz:   2.0,
		Ram:         &psm.Memory{Value: 8, Unit: psm.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
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
	// testSearchLaptop(laptopClient)
	// testUploadImage(laptopClient)
	testRateLaptop(laptopClient)

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

func rateLaptop(laptopClient psm.LaptopServiceClient, laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.RateLaptop(ctx)
	if err != nil {
		return err
	}
	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("Cannot receive response: %v", err)
				return
			}

			log.Printf("Laptop with ID %s has average score: %f", res.GetLaptopId(), res.GetAverageScore())
		}
	}()

	for i, laptopID := range laptopIDs {
		req := &psm.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}

		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("Cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}
		log.Print("send request", req)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("Cannot close stream: %v", err)
	}

	err = <-waitResponse
	return err
}

func testRateLaptop(laptopClient psm.LaptopServiceClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.Laptop()
		CreateRandomLaptop(laptopClient, laptop)
		laptopIDs[i] = laptop.GetId()
	}

	scores := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Print("rate laptop (y/n)?")
		var answer string
		fmt.Scan(&answer)
		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := rateLaptop(laptopClient, laptopIDs, scores)
		if err != nil {
			log.Fatal("Cannot rate laptop", err)
		}

	}

}
