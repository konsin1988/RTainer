package grpcserver

import (
    "context"

    pb "konsin1988/agent/proto"
)

func (s *Server) DockerInfo(
    ctx context.Context,
    req *pb.DockerInfoRequest,
) (*pb.DockerInfoResponse, error) {

    info, err := s.containerSvc.DockerInfo(ctx)
    if err != nil {
        return nil, err
    }

    return &pb.DockerInfoResponse{
        ServerVersion:     info.ServerVersion,
        ApiVersion:        info.APIVersion,
        OperatingSystem:   info.OperatingSystem,
        Architecture:      info.Architecture,
        KernelVersion:     info.KernelVersion,
        Containers:        info.Containers,
        ContainersRunning: info.ContainersRunning,
        ContainersPaused:  info.ContainersPaused,
        ContainersStopped: info.ContainersStopped,
        Images:            info.Images,
        Volumes:           info.Volumes,
        Driver:            info.Driver,
        MemTotal:          info.MemTotal,
        Ncpu:              info.NCPU,
        DockerRootDir:     info.DockerRootDir,
        Name:              info.Name,
        Id:                info.ID,
    }, nil
}
