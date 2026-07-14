package service

import (
	"context"

	"konsin1988/agent/docker"
)

type VolumeService struct {
    docker *docker.Client
}

func NewVolumeService(d *docker.Client) *VolumeService {
    return &VolumeService{docker: d}
}


// ---------------- LIST VOLUMES -----------------------------
func (s *VolumeService) ListVolumes(
    ctx context.Context,
) ([]docker.Volume, error) {

    return s.docker.ListVolumes(ctx)
}
