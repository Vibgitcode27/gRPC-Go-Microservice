package serializer

import (
	"fmt"
	"os"

	// "io/ioutil"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// This function writes protocal buffer message to binary file
func WriteProtobufToBinayFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	// err = ioutil.WriteFile(filename, data, 0644)   // ioutil.WriteFile is deprecated
	// if err != nil {
	// 	return fmt.Errorf("cannot write binary data to file: %w", err)
	// }

	// Create or open the file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	// Write the binary data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	// data, err := ioutil.ReadFile(filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary data: %w", err)
	}

	return nil
}

// WriteProtobufToJSONFile writes protocol buffer message to JSON file
func WriteProtobufToJSONFile(messageOk proto.Message, filename string) error {
	// Marshal the message to JSON format
	data, err := ProtobufToJson(messageOk)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}

	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON data to file: %w", err)
	}

	return nil
}

func ProtobufToJson(message proto.Message) (string, error) {

	protoMsg := message.ProtoReflect().Interface()

	// Convert protoMsg to protoiface.MessageV1
	protoMessage, ok := protoMsg.(protoiface.MessageV1)
	if !ok {
		return "", fmt.Errorf("cannot convert proto message to protoiface.MessageV1")
	}

	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	json, err := marshaler.MarshalToString(protoMessage)
	if err != nil {
		return "", err
	}
	return json, nil
}
