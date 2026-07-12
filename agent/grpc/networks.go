package grpcserver

import (
    "context"

    pb "konsin1988/agent/proto"
)

// ------------------------- LIST NETWORKS -----------------
func (s *Server) ListNetworks(
    ctx context.Context,
    req *pb.ListNetworksRequest,
) (*pb.ListNetworksResponse, error) {


    networks, err := s.networkSvc.ListNetworks(ctx)

    if err != nil {
        return nil, err
    }


    resp := &pb.ListNetworksResponse{}


    for _, n := range networks {

        item := &pb.Network{
            Id:     n.ID,
            Name:   n.Name,
            Driver: n.Driver,
            Scope:  n.Scope,
        }

				for _, c := range n.Containers {
            item.Containers = append(
                item.Containers,
                &pb.NetContainer{
                    Id:          c.ID,
                    Name:        c.Name,
                    Ipv4Address: c.IPv4Address,
                },
            )
        }

        resp.Networks = append(
            resp.Networks,
            item,
        )
    }


    return resp, nil
}


// -------------------- CREATE NETWORK ------------------
func (s *Server) CreateNetwork(
    ctx context.Context,
    req *pb.CreateNetworkRequest,
) (*pb.CreateNetworkResponse, error) {


    id, err := s.networkSvc.CreateNetwork(
        ctx,
        req,
    )

    if err != nil {
        return nil, err
    }


    return &pb.CreateNetworkResponse{
        Id: id,
    }, nil
}

// ------------------ REMOVE NETWORK -----------------------
func (s *Server) RemoveNetwork(
    ctx context.Context,
    req *pb.RemoveNetworkRequest,
) (*pb.RemoveNetworkResponse, error) {


    err := s.networkSvc.RemoveNetwork(
        ctx,
        req,
    )

    if err != nil {
        return nil, err
    }


    return &pb.RemoveNetworkResponse{}, nil
}
