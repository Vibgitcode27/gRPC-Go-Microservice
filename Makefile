gen :
	protoc --go_out=. --go-grpc_out=. proto/process_manager.proto proto/memory_message.proto

clean :
	rm -rf psm

run :
	go build
	./grpc