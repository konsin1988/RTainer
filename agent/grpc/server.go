package grpcserver

import (

    pb "konsin1988/agent/proto"
    "konsin1988/agent/service"

    "google.golang.org/grpc"
)

type Server struct {
    pb.UnimplementedAgentServiceServer
    pb.UnimplementedContainerServiceServer
    pb.UnimplementedImageServiceServer
    pb.UnimplementedNetworkServiceServer
    pb.UnimplementedVolumeServiceServer
    containerSvc *service.ContainerService
		imageSvc *service.ImageService
		networkSvc *service.NetworkService
		volumeSvc *service.VolumeService
}

func New(
	container *service.ContainerService, 
	image *service.ImageService, 
	network *service.NetworkService,
	volume *service.VolumeService,
) *Server {
	return &Server{containerSvc: container, imageSvc: image, networkSvc: network, volumeSvc: volume}
}

func (s *Server) Register(grpcSrv *grpc.Server) {
    pb.RegisterAgentServiceServer(grpcSrv, s)
    pb.RegisterContainerServiceServer(grpcSrv, s)
    pb.RegisterImageServiceServer(grpcSrv, s)
    pb.RegisterNetworkServiceServer(grpcSrv, s)
    pb.RegisterVolumeServiceServer(grpcSrv, s)

}




