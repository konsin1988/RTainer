package grpcserver

import (
    "context"

    pb "konsin1988/agent/proto"

		"google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// -------------------------------------------------- LIST NETWORKS 
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

        item := &pb.NetworkListItem{
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

// ------------------------------------------------------- INSPECT NETWORK 
func (s *Server) InspectNetwork(
    ctx context.Context,
    req *pb.NetworkRequest,
) (*pb.NetworkResponse, error) {

    resp, err := s.networkSvc.InspectNetwork(ctx, req.Id)
    if err != nil {
        return nil, err
    }
		return &pb.NetworkResponse{Network: resp}, nil
}

// ------------------------------------------------------- CREATE NETWORK
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

// -------------------------------------------------------- REMOVE NETWORK
func (s *Server) RemoveNetwork(
    ctx context.Context,
    req *pb.NetworkRequest,
) (*pb.NetworkResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }
    resp, err := s.networkSvc.InspectNetwork(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    err = s.networkSvc.RemoveNetwork(
        ctx,
        req.Id,
    )
    if err != nil {
        return nil, err
    }

		return &pb.NetworkResponse{Network: resp}, nil
}


// ------------------------------------------------------ CONNECT NETWORK
func (s *Server) ConnectNetwork(
    ctx context.Context,
    req *pb.ConnectNetworkRequest,
) (*pb.NetworkResponse, error) {

		if req.GetNetworkId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "network_id is required",
        )
    }

    if req.GetContainerId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "container_id is required",
        )
    }
		
    resp, err := s.networkSvc.InspectNetwork(ctx, req.NetworkId)
    if err != nil {
        return nil, err
    }

    err = s.networkSvc.ConnectNetwork(
        ctx,
        req.NetworkId,
				req.ContainerId,
				req.Endpoint,
    )

    if err != nil {
        return nil, err
    }

		return &pb.NetworkResponse{Network: resp}, nil
}

// ------------------------------------------------------ DISCONNECT NETWORK
func (s *Server) DisconnectNetwork(
    ctx context.Context,
    req *pb.DisconnectNetworkRequest,
) (*pb.NetworkResponse, error) {

		if req.GetNetworkId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "network_id is required",
        )
    }
    if req.GetContainerId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "container_id is required",
        )
    }
		

		err := s.networkSvc.DisconnectNetwork(
        ctx,
        req.NetworkId,
				req.ContainerId,
				req.Forse,
    )

    if err != nil {
        return nil, err
    }

    resp, err := s.networkSvc.InspectNetwork(ctx, req.NetworkId)
    if err != nil {
        return nil, err
    }

		return &pb.NetworkResponse{Network: resp}, nil
}
