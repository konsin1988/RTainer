package service

import (
	"context"

	"konsin1988/agent/docker"
	"konsin1988/agent/proto"

	"github.com/docker/docker/api/types/filters"
	"google.golang.org/protobuf/types/known/structpb"
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

// -------------------------------------------- INSPECT VOLUME
func (s *VolumeService) InspectVolume(
		ctx context.Context,
		name string,
) (*proto.InspectVolumeResponse, error){
		resp, err := s.docker.InspectVolume(ctx, name)
		if err != nil{
			return nil, err
		}

		statusStruct, err := structpb.NewStruct(resp.Status)

		return &proto.InspectVolumeResponse{
			Volume: &proto.VolumeInfo{
				Name:  				resp.Name,   
				Driver: 			resp.Driver, 
				MountPoint:		resp.Mountpoint, 
				CreatedAt: 		resp.CreatedAt, 
				Status:   		statusStruct, 
				Scope:   			resp.Scope, 
				Labels: 			resp.Labels, 
				Options: 			resp.Options, 
			},
		}, nil
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

// --------------------------------------------------- PRUNE VOLUME 
func (s *VolumeService) PruneVolume(
	ctx context.Context,
	req *proto.PruneVolumeRequest,
)(*proto.PruneVolumeResponse, error){

		pruneFilters := volumePruneFiltersFromProto(req)

    report, err := s.docker.PruneVolume(ctx, pruneFilters)
    if err != nil {
        return nil, err 
    }

    return &proto.PruneVolumeResponse{
        DeletedVolumeIds: report.VolumesDeleted,
    }, nil
}

func volumePruneFiltersFromProto(
    req *proto.PruneVolumeRequest,
) filters.Args {
    f := filters.NewArgs()

    for key, value := range req.GetFilters() {
        f.Add(key, value)
    }

    return f
}
