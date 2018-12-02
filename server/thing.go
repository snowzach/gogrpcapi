package server

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/snowzach/gogrpcapi/gogrpcapi"
	"github.com/snowzach/gogrpcapi/server/rpc"
	"github.com/snowzach/gogrpcapi/store"
)

// ThingFind returns all things
func (s *Server) ThingFind(ctx context.Context, _ *emptypb.Empty) (*rpc.ThingFindResponse, error) {

	bs, err := s.thingStore.ThingFind(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.ThingFindResponse{
		Data: bs,
	}, nil

}

// ThingGet fetches a thing by ID
func (s *Server) ThingGet(ctx context.Context, request *rpc.ThingId) (*gogrpcapi.Thing, error) {

	if request.Id == "" {
		return nil, grpc.Errorf(codes.Internal, "Invalid ID")
	}
	b, err := s.thingStore.ThingGetById(ctx, request.Id)
	if err == store.ErrNotFound {
		return nil, grpc.Errorf(codes.NotFound, "Not Found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return b, nil

}

// ThingSave creates or updates a thing
func (s *Server) ThingSave(ctx context.Context, b *gogrpcapi.Thing) (*rpc.ThingId, error) {

	thingID, err := s.thingStore.ThingSave(ctx, b)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &rpc.ThingId{
		Id: thingID,
	}, nil

}

// ThingDelete deletes a thing
func (s *Server) ThingDelete(ctx context.Context, request *rpc.ThingId) (*emptypb.Empty, error) {

	if request.Id == "" {
		return nil, grpc.Errorf(codes.Internal, "Invalid ID")
	}
	err := s.thingStore.ThingDeleteById(ctx, request.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &emptypb.Empty{}, nil

}
