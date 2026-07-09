package grpcserver

import (
    "context"
    "google.golang.org/grpc"
		"strings"

    pb "konsin1988/agent/proto"
    "konsin1988/agent/service"
)

type Server struct {
    pb.UnimplementedAgentServiceServer
    containerSvc *service.ContainerService
		imageSvc *service.ImageService
}

func New(container *service.ContainerService, image *service.ImageService) *Server {
		return &Server{containerSvc: container, imageSvc: image}
}

func (s *Server) Register(grpcSrv *grpc.Server) {
    pb.RegisterAgentServiceServer(grpcSrv, s)
}


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

// -----------------------------------
// LIST IMAGES
// ------------------------------------
func (s *Server) ListImages(
		ctx context.Context,
		req *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
    images, err := s.imageSvc.ListImages(ctx)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListImagesResponse{}

    for _, c := range images {
        resp.Images = append(resp.Images, &pb.Image{
            Id:     	c.ID,
						RepoTags: strings.Join(c.RepoTags, ","), 
						Size:			c.Size,     
						Created: 	c.Created,   
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
