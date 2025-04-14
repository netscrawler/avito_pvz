package pvzgrpc

import (
	"context"
	"errors"

	"avito_pvz/internal/models"
	"avito_pvz/internal/models/domain"

	pvzv1 "github.com/netscrawler/pvz_proto/gen/go/pvz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZ interface {
	GetAllPVZ(ctx context.Context) (domain.PVZList, error)
}

type serverAPI struct {
	pvzv1.UnimplementedPVZServiceServer
	pvz PVZ
}

//nolint:exhaustruct
func Register(grpcServer *grpc.Server, pvz PVZ) {
	pvzv1.RegisterPVZServiceServer(grpcServer, &serverAPI{pvz: pvz})
}

func (s *serverAPI) GetPVZList(
	ctx context.Context,
	in *pvzv1.GetPVZListRequest,
) (*pvzv1.GetPVZListResponse, error) {
	list, err := s.pvz.GetAllPVZ(ctx)
	if errors.Is(err, models.ErrPVZNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := &pvzv1.GetPVZListResponse{
		Pvzs: []*pvzv1.PVZ{},
	}

	for _, v := range list {
		p := pvzv1.PVZ{
			Id:               v.ID.String(),
			RegistrationDate: timestamppb.New(v.RegistrationDate),
			City:             string(v.City),
		}
		out.Pvzs = append(out.Pvzs, &p)
	}

	return out, nil
}
