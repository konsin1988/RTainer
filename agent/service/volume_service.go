package service

import (
	"context"

	"konsin1988/agent/docker"
	"konsin1988/agent/proto"
)

type VolumeService struct {
    docker *docker.Client
}

func NewVolumeService(d *docker.Client) *VolumeService {
    return &VolumeService{docker: d}
}


// ------------------------------------------ LIST VOLUMES
func (s *VolumeService) ListVolumes(
    ctx context.Context,
) ([]docker.Volume, error) {

    return s.docker.ListVolumes(ctx)
}

// ------------------------------------------- CREATE VOLUME
func (s *VolumeService) CreateVolume(
    ctx context.Context,
    req *proto.CreateVolumeRequest,
) error {

    return s.docker.CreateVolume(
        ctx,
        docker.CreateVolumeRequest{
            Name:    req.Name,
            Driver:  req.Driver,
            Labels:  req.Labels,
            Options: req.Options,
        },
    )
}

// -------------------------------------------- REMOVE VOLUME
func (s *VolumeService) RemoveVolume(
    ctx context.Context,
    req *proto.RemoveVolumeRequest,
) error {

    return s.docker.RemoveVolume(
        ctx,
        req.Name,
        req.Force,
    )
}
