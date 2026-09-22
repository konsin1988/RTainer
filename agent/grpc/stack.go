package grpcserver

import (
    "context"

    pb "konsin1988/agent/proto"

		_ "google.golang.org/grpc/codes"
    _ "google.golang.org/grpc/status"
)

// -------------------------------------------------- LIST NETWORKS 
func (s *Server) ListStacks(
    ctx context.Context,
    req *pb.ListStackRequest,
) (*pb.ListStackResponse, error) {

  return s.stackSvc.ListStacks(ctx)
}
