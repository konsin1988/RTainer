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

// ----------------------- LIST NETWORKS -------------------
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


// -------------------- CREATE NETWORK --------------------
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

// --------------------- REMOVE NETWORK --------------------
func (c *Client) RemoveNetwork(
    ctx context.Context,
    id string,
) error {
    return c.cli.NetworkRemove(
        ctx,
        id,
    )
}
