package grpcserver

import (

    pb "konsin1988/agent/proto"
    "konsin1988/agent/service"

    "google.golang.org/grpc"
)

type Server struct {
    pb.UnimplementedAgentServiceServer
    containerSvc *service.ContainerService
		imageSvc *service.ImageService
		networkSvc *service.NetworkService
}

func New(
	container *service.ContainerService, 
	image *service.ImageService, 
	network *service.NetworkService,
) *Server {
	return &Server{containerSvc: container, imageSvc: image, networkSvc: network}
}

func (s *Server) Register(grpcSrv *grpc.Server) {
    pb.RegisterAgentServiceServer(grpcSrv, s)
}




