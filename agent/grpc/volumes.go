package grpcserver

import (
	"context"
    
	pb "konsin1988/agent/proto"

	"google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

// -------------------------------------------- LIST VOLUMES
func (s *Server) ListVolumes(
    ctx context.Context,
    req *pb.ListVolumesRequest,
) (*pb.ListVolumesResponse, error) {

    return s.volumeSvc.ListVolumes(ctx)
}

// ------------------------------------------------- CREATE VOLUME
func (s *Server) CreateVolume(
    ctx context.Context,
    req *pb.CreateVolumeRequest,
) (*pb.VolumeResponse, error) {

    err := s.volumeSvc.CreateVolume(ctx, req)
    if err != nil {
        return nil, err
    }

    return &pb.VolumeResponse{}, nil
}


// ------------------------------------------------- REMOVE VOLUME 
func (s *Server) RemoveVolume(
    ctx context.Context,
    req *pb.RemoveVolumeRequest,
) (*pb.VolumeResponse, error) {
    err := s.volumeSvc.RemoveVolume(ctx, req)
    if err != nil {
        return nil, err
    }

    return &pb.VolumeResponse{}, nil
}

// ------------------------------------------------ INSPECT VOLUME
func (s *Server) InspectVolume(
		ctx context.Context,
		req *pb.InspectVolumeRequest,
) (*pb.InspectVolumeResponse, error) {
		if req.GetName() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }
		resp, err := s.volumeSvc.InspectVolume(ctx, req.Name)

		if err != nil {
			return nil, err
		}
		return resp, nil
}


// ----------------------------------------------------------------- PRUNE VOLUME 
func (s *Server) PruneVolume(
	ctx context.Context,
	req *pb.PruneVolumeRequest,
) (*pb.PruneVolumeResponse, error) {

	return s.volumeSvc.PruneVolume(ctx, req) 
}
