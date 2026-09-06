package docker

import (
    "context"

    "github.com/docker/docker/api/types/network"
)

type NetContainer struct {
		ID 						string
		Name					string
		IPv4Address 	string
}

type Network struct {
    ID         string
    Name       string
    Driver     string
    Scope      string
    Containers []NetContainer
}

type IPAMConfig struct {
    Subnet  string
    Gateway string
    IPRange string
}

type NetworkInfo struct {
    ID         string
    Name       string
    Driver     string
    Scope      string

    Internal   bool
    Attachable bool
    Ingress    bool

    IPAM       []IPAMConfig

    Options    map[string]string
    Labels     map[string]string

    Containers []NetContainer
}

// ------------------------------------------------- LIST NETWORKS 
func (c *Client) ListNetworks(
    ctx context.Context,
) ([]Network, error) {

    networks, err := c.cli.NetworkList(
        ctx,
        network.ListOptions{},
    )

    if err != nil {
        return nil, err
    }


    result := make([]Network, 0, len(networks))


    for _, n := range networks {
				inspect, err := c.cli.NetworkInspect(
        		ctx,
        		n.ID,
        		network.InspectOptions{},
    		)
    		if err != nil {
    		    return nil, err
    		}

    		item := Network{
    		    ID:     inspect.ID,
    		    Name:   inspect.Name,
    		    Driver: inspect.Driver,
    		    Scope:  inspect.Scope,
    		}

				for id, endpoint := range inspect.Containers {
				
				    item.Containers = append(
				        item.Containers,
				        NetContainer{
				            ID:          id,
				            Name:        endpoint.Name,
				            IPv4Address: endpoint.IPv4Address,
				        },
				    )
				}

        result = append(result, item)
    }
    return result, nil
}


// ------------------------------------------------------ CREATE NETWORK 
func (c *Client) CreateNetwork(
    ctx context.Context,
    name string,
    driver string,
    internal bool,
) (string, error) {


    resp, err := c.cli.NetworkCreate(
        ctx,
        name,
        network.CreateOptions{
            Driver:   driver,
            Internal: internal,
        },
    )

    if err != nil {
        return "", err
    }


    return resp.ID, nil
}

// ------------------------------------------------------ REMOVE NETWORK 
func (c *Client) RemoveNetwork(
    ctx context.Context,
    id string,
) error {
    return c.cli.NetworkRemove(
        ctx,
        id,
    )
}


// ------------------------------------------------------ INSPECT NETWORK 
func (c *Client) InspectNetwork(
    ctx context.Context,
    id string,
) (NetworkInfo, error) {

    n, err := c.cli.NetworkInspect(
        ctx,
        id,
        network.InspectOptions{},
    )
    if err != nil {
        return NetworkInfo{}, err
    }

    info := NetworkInfo{
        ID:         n.ID,
        Name:       n.Name,
        Driver:     n.Driver,
        Scope:      n.Scope,
        Internal:   n.Internal,
        Attachable: n.Attachable,
        Ingress:    n.Ingress,
        Options:    n.Options,
        Labels:     n.Labels,
    }

    for _, cfg := range n.IPAM.Config {
        info.IPAM = append(info.IPAM, IPAMConfig{
            Subnet:  cfg.Subnet,
            Gateway: cfg.Gateway,
            IPRange: cfg.IPRange,
        })
    }

    for id, endpoint := range n.Containers {
        info.Containers = append(info.Containers, NetContainer{
            ID:          id,
            Name:        endpoint.Name,
            IPv4Address: endpoint.IPv4Address,
        })
    }

    return info, nil
}

// ------------------------------------------------------ CONNECT NETWORK 
func (c *Client) ConnectNetwork(
    ctx context.Context,
    networkID string,
		containerID string,
		endpoint *network.EndpointSettings,
) error {
		return c.cli.NetworkConnect(
        ctx,
        networkID,
        containerID,
        endpoint,
    )
}

// ------------------------------------------------------ DISCONNECT NETWORK 
func (c *Client) DisconnectNetwork(
    ctx context.Context,
    networkID string,
		containerID string,
		forse bool,
) error {
		return c.cli.NetworkDisconnect(
        ctx,
        networkID,
        containerID,
				forse,
    )
}
