package serializer_test

import (
	"grpc/psm"
	"grpc/sample"
	"grpc/serializer"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../temp/laptop.bin"

	laptop := sample.Laptop()
	err := serializer.WriteProtobufToBinayFile(laptop, binaryFile)
	require.NoError(t, err)

	laptop2 := &psm.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	err = serializer.WriteProtobufToJSONFile(laptop, "../temp/laptop.json")
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop, laptop2), "Test Successful! ")
}
