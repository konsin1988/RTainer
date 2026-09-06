package service


import (
    "context"

    "konsin1988/agent/docker"
    "konsin1988/agent/proto"
    "github.com/docker/docker/api/types/network"
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

// ------------------------------------------------------ LIST NETWORKS
func (s *NetworkService) ListNetworks(
    ctx context.Context,
) ([]docker.Network, error) {

    return s.docker.ListNetworks(ctx)
}


// ------------------------------------------------------- CREATE NETWORK 
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


// ------------------------------------------------------- REMOVE NETWORK
func (s *NetworkService) RemoveNetwork(
    ctx context.Context,
    id string,
) error {

    return s.docker.RemoveNetwork(
        ctx,
        id,
    )
}


// ----------------------------------------------------- INSPECT NETWORK 
func (s *NetworkService) InspectNetwork(
    ctx context.Context,
    id string,
) (*proto.NetworkInfo, error) {

    n, err := s.docker.InspectNetwork(ctx, id)
    if err != nil {
        return nil, err
    }

    resp := &proto.NetworkInfo {
        Id:         n.ID,
        Name:       n.Name,
        Driver:     n.Driver,
        Scope:      n.Scope,
        Internal:   n.Internal,
        Attachable: n.Attachable,
        Ingress:    n.Ingress,
        Labels:     n.Labels,
    }

    for _, cfg := range n.IPAM {
        resp.Ipam = append(resp.Ipam, &proto.IPAMConfig{
            Subnet:  cfg.Subnet,
            Gateway: cfg.Gateway,
            IpRange: cfg.IPRange,
        })
    }

    for _, c := range n.Containers {
        resp.Containers = append(resp.Containers, &proto.NetContainer{
            Id:          c.ID,
            Name:        c.Name,
            Ipv4Address: c.IPv4Address,
        })
    }
    return resp, nil 
}

// ----------------------------------------------------- CONNECT NETWORK 
func (s *NetworkService) ConnectNetwork(
    ctx context.Context,
    network_id string,
		container_id string,
		endpoint *proto.EndpointSettings,
) (error) {

		dockerEndpoint := endpointSettingsFromProto(endpoint)

    return s.docker.ConnectNetwork(ctx, network_id, container_id, dockerEndpoint )
}

func endpointSettingsFromProto(
    endpoint *proto.EndpointSettings,
) *network.EndpointSettings {

    if endpoint == nil {
        return nil
    }

    return &network.EndpointSettings{
        IPAMConfig: &network.EndpointIPAMConfig{
            IPv4Address: endpoint.GetIpv4Address(),
            IPv6Address: endpoint.GetIpv6Address(),
        },

        MacAddress: endpoint.GetMacAddress(),

        Aliases: endpoint.GetAliases(),
    }
}

// ----------------------------------------------------- DISCONNECT NETWORK 
func (s *NetworkService) DisconnectNetwork(
    ctx context.Context,
    network_id string,
		container_id string,
		forse bool,
) (error) {

    return s.docker.DisconnectNetwork(ctx, network_id, container_id, forse)
}
