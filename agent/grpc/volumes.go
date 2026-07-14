package grpcserver

import (
	"context"
    
	pb "konsin1988/agent/proto"
)

// -------------------------------------------- LIST VOLUMES
func (s *Server) ListVolumes(
    ctx context.Context,
    req *pb.ListVolumesRequest,
) (*pb.ListVolumesResponse, error) {

    volumes, err := s.volumeSvc.ListVolumes(ctx)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListVolumesResponse{}

    for _, v := range volumes {
        resp.Volumes = append(resp.Volumes, &pb.Volume{
            Name:       v.Name,
            Driver:     v.Driver,
            Mountpoint: v.Mountpoint,
            Labels:     v.Labels,
            Scope:      v.Scope,
        })
    }

    return resp, nil
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
