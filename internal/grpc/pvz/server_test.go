package pvzgrpc_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	pvzgrpc "avito_pvz/internal/grpc/pvz"
	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"

	"github.com/google/uuid"
	pvzv1 "github.com/netscrawler/pvz_proto/gen/go/pvz"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestGetPVZList(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mockPVZ *pvzgrpc.MockPVZ) // Параметр для настройки моков
		wantErr    error
		wantCode   codes.Code
	}{
		{
			name: "successful_response",
			setupMocks: func(mockPVZ *pvzgrpc.MockPVZ) {
				mockPVZ.On("GetAllPVZ", mock.Anything).Return(domain.PVZList{
					{
						ID:               (*domain.PVZID)(&uuid.Max),
						RegistrationDate: time.Time{},
						City:             "Москва",
					},
					{
						ID:               (*domain.PVZID)(&uuid.Max),
						RegistrationDate: time.Time{},
						City:             "Казань",
					},
				}, nil)
			},
			wantErr:  nil,
			wantCode: codes.OK,
		},
		{
			name: "pvz_not_found",
			setupMocks: func(mockPVZ *pvzgrpc.MockPVZ) {
				mockPVZ.On("GetAllPVZ", mock.Anything).Return(nil, models.ErrPVZNotFound)
			},
			wantErr:  models.ErrPVZNotFound,
			wantCode: codes.NotFound,
		},
		{
			name: "internal_error",
			setupMocks: func(mockPVZ *pvzgrpc.MockPVZ) {
				mockPVZ.On("GetAllPVZ", mock.Anything).Return(nil, errors.New("internal error"))
			},
			wantErr:  models.ErrInternal,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZ := new(pvzgrpc.MockPVZ)
			tt.setupMocks(mockPVZ)

			lis, err := net.Listen("tcp", ":0") // Используем случайный порт
			require.NoError(t, err)
			defer lis.Close()

			grpcServer := grpc.NewServer()
			pvzgrpc.Register(grpcServer, mockPVZ)

			ready := make(chan struct{})
			done := make(chan struct{})

			go func() {
				close(ready)
				err := grpcServer.Serve(lis)
				require.NoError(t, err)
				close(done)
			}()

			<-ready

			conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			client := pvzv1.NewPVZServiceClient(conn)

			resp, err := client.GetPVZList(context.Background(), &pvzv1.GetPVZListRequest{})
			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			mockPVZ.AssertExpectations(t)

			grpcServer.GracefulStop()
			<-done
		})
	}
}
