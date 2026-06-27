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
    req *pb.StopContainerRequest,
) (*pb.StopContainerResponse, error) {

    err := s.containerSvc.StopContainer(ctx, req.Id, int(req.Timeout))
    if err != nil {
        return &pb.StopContainerResponse{
            Success: false,
        }, err
    }

    return &pb.StopContainerResponse{
        Success: true,
    }, nil
}
