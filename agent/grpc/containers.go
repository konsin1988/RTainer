package grpcserver

import (
    "context"

    pb "konsin1988/agent/proto"
)


// ---------------------------------------
// LIST CONTAINERS 
// --------------------------------------
func (s *Server) ListContainers(
    ctx context.Context,
    req *pb.ListContainersRequest,
) (*pb.ListContainersResponse, error) {

    ctrs, err := s.containerSvc.ListContainers(ctx)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListContainersResponse{}

    for _, c := range ctrs {
        resp.Containers = append(resp.Containers, &pb.Container{
            Id:     c.ID,
            Name:   c.Name,
            Image:  c.Image,
            Status: c.Status,
        })
    }

    return resp, nil
}

// ------------------------------------
// STOP CONTAINER
// ------------------------------------
func (s *Server) StopContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.StopContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}


// ------------------------------------
// START CONTAINER
// ------------------------------------
func (s *Server) StartContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.StartContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}

// ------------------------------------
// REMOVE CONTAINER
// ------------------------------------
func (s *Server) RemoveContainer(
    ctx context.Context,
    req *pb.RemoveContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.RemoveContainer(ctx, req.Id, req.Force, req.RemoveVolumes)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}


// ------------------------------------
// RUN CONTAINER
// ------------------------------------
func (s *Server) RunContainer(
    ctx context.Context,
    req *pb.RunContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.RunContainer(ctx, req)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}



