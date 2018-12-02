package server

import (
	"net/http"

	"github.com/snowzach/gogrpcapi/server/rpc"
)

// SetupRoutes configures all the routes for this service
func (s *Server) SetupRoutes() {

	// Register our routes - you need at aleast one route
	s.router.Get("/none", func(w http.ResponseWriter, r *http.Request) {})

	// Register RPC Services
	rpc.RegisterVersionRPCServer(s.grpcServer, s)
	s.gwReg(rpc.RegisterVersionRPCHandlerFromEndpoint)
	rpc.RegisterThingRPCServer(s.grpcServer, s)
	s.gwReg(rpc.RegisterThingRPCHandlerFromEndpoint)

}
