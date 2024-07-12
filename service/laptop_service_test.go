package service_test

import (
	"context"
	"grpc/psm"
	"grpc/sample"
	"grpc/service"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	t.Parallel()

	laptopNoId := sample.Laptop()
	laptopNoId.Id = ""

	laptopInvalidId := sample.Laptop()
	laptopInvalidId.Id = "invalid-uuid"

	loadDuplicateId := sample.Laptop()
	storeDuplicateId := service.NewInMemoryLaptopStore()

	err := storeDuplicateId.Save(loadDuplicateId)
	if err != nil {
		t.Errorf("cannot save laptop to store: %v", err)
	}

	testCases := []struct {
		name        string
		laptop      *psm.Laptop
		laptopStore service.LaptopStore
		code        codes.Code
	}{
		{
			name:        "success_with_id",
			laptop:      sample.Laptop(),
			laptopStore: service.NewInMemoryLaptopStore(),

			code: codes.OK,
		},
		{
			name:        "success_no_id",
			laptop:      laptopNoId,
			laptopStore: service.NewInMemoryLaptopStore(), code: codes.OK,
		},
		{
			name:        "failure_invalid_id",
			laptop:      laptopInvalidId,
			laptopStore: service.NewInMemoryLaptopStore(),
			code:        codes.InvalidArgument,
		},
		{
			name:        "failure_duplicate_id",
			laptop:      loadDuplicateId,
			laptopStore: storeDuplicateId,
			code:        codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &psm.CreateLaptopRequest{
				Laptop: tc.laptop,
			}

			imageStore := service.NewDiskImageStore("tmp")
			ratingStore := service.NewInMemoryRatingStore()

			server := service.NewLaptopService(tc.laptopStore, imageStore, ratingStore)
			res, err := server.CreateLaptop(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)

				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}
}
