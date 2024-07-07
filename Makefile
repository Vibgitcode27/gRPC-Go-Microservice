gen :
	protoc --go_out=. --go-grpc_out=. proto/process_manager.proto proto/memory_message.proto proto/keyboard_message.proto proto/laptop_message.proto proto/screen_message.proto proto/storage_message.proto

clean :
	rm -rf psm

run :
	go build
	./grpc