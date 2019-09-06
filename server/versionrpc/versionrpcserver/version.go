package versionrpcserver

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"

	"github.com/snowzach/gogrpcapi/conf"
	"github.com/snowzach/gogrpcapi/server/versionrpc"
)

type versionRPCServer struct{}

// New returns a new version server
func New() versionrpc.VersionRPCServer {
	return versionRPCServer{}
}

// Version returns the version
func (vs versionRPCServer) Version(ctx context.Context, _ *emptypb.Empty) (*versionrpc.VersionResponse, error) {

	return &versionrpc.VersionResponse{
		Version: conf.GitVersion,
	}, nil

}
