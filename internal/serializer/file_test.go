package serializer

import (
	"testing"

	"github.com/caiofernandes00/playing-with-golang/grpc/pkg/proto/pb"
	"github.com/caiofernandes00/playing-with-golang/grpc/internal/sample"
	"github.com/stretchr/testify/require"
)

func TestFileSerialzer(t *testing.T) {
	t.Parallel()

	binaryFile := "../../tmp/laptop.bin"
	jsonFile := "../../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)

	err = WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)
}
