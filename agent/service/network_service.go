package service


import (
    "context"

    "konsin1988/agent/docker"
    "konsin1988/agent/proto"
)


type NetworkService struct {
    docker *docker.Client
}


func NewNetworkService(
    d *docker.Client,
) *NetworkService {

    return &NetworkService{
        docker: d,
    }
}

// ---------------------LIST NETWORKS ------------------
func (s *NetworkService) ListNetworks(
    ctx context.Context,
) ([]docker.Network, error) {

    return s.docker.ListNetworks(ctx)
}


// ------------------------ CREATE NETWORK -------------
func (s *NetworkService) CreateNetwork(
    ctx context.Context,
    req *proto.CreateNetworkRequest,
) (string, error) {

    return s.docker.CreateNetwork(
        ctx,
        req.Name,
        req.Driver,
        req.Internal,
    )
}


// ------------------------ REMOVE NETWORK -------------------
func (s *NetworkService) RemoveNetwork(
    ctx context.Context,
    req *proto.RemoveNetworkRequest,
) error {

    return s.docker.RemoveNetwork(
        ctx,
        req.Id,
    )
}
