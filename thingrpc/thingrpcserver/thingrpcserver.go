package thingrpcserver

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/snowzach/gogrpcapi/store"
	"github.com/snowzach/gogrpcapi/thingrpc"
)

type thingRPCServer struct {
	thingStore thingrpc.ThingStore
}

// New returns a new rpc server
func New(ts thingrpc.ThingStore) (thingrpc.ThingRPCServer, error) {

	return newServer(ts)

}

func newServer(ts thingrpc.ThingStore) (*thingRPCServer, error) {

	return &thingRPCServer{
		thingStore: ts,
	}, nil

}

// AuthFuncOverride is used if you want to override default authentication for any endpoint
// This disables all authentication for any thingRPC calls
func (s *thingRPCServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

// ThingFind returns all things
func (s *thingRPCServer) ThingFind(ctx context.Context, _ *emptypb.Empty) (*thingrpc.ThingFindResponse, error) {

	bs, err := s.thingStore.ThingFind(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &thingrpc.ThingFindResponse{
		Data: bs,
	}, nil

}

// ThingGet fetches a thing by ID
func (s *thingRPCServer) ThingGet(ctx context.Context, request *thingrpc.ThingId) (*thingrpc.Thing, error) {

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
func (s *thingRPCServer) ThingSave(ctx context.Context, b *thingrpc.Thing) (*thingrpc.ThingId, error) {

	thingID, err := s.thingStore.ThingSave(ctx, b)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &thingrpc.ThingId{
		Id: thingID,
	}, nil

}

// ThingDelete deletes a thing
func (s *thingRPCServer) ThingDelete(ctx context.Context, request *thingrpc.ThingId) (*emptypb.Empty, error) {

	if request.Id == "" {
		return nil, grpc.Errorf(codes.Internal, "Invalid ID")
	}
	err := s.thingStore.ThingDeleteById(ctx, request.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%s", err)
	}

	return &emptypb.Empty{}, nil

}
