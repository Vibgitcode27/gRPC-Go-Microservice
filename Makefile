gen :
	protoc --go_out=. --go-grpc_out=. proto/process_manager.proto proto/memory_message.proto proto/keyboard_message.proto proto/laptop_message.proto proto/screen_message.proto proto/storage_message.proto proto/laptop_service.proto

clean :
	rm -rf psm
 
# run :
# 	go build
# 	./grpc

server : 
	go run cmd/server/main.go -port 8080

client :
	go run cmd/client/main.go -address 0.0.0.0:8080
test:
	rm -rf temp/*
	go test -v ./...