package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
)

type Volume struct {
    Name       string
    Driver     string
    Mountpoint string
    Labels     map[string]string
    Scope      string
}

type CreateVolumeRequest struct {
    Name    string
    Driver  string
    Labels  map[string]string
    Options map[string]string
}


// ---------------------------------------- LIST VOLUMES 
func (c *Client) ListVolumes(
    ctx context.Context,
) ([]Volume, error) {

    resp, err := c.cli.VolumeList(
        ctx,
        volume.ListOptions{},
    )
    if err != nil {
        return nil, err
    }

    result := make([]Volume, 0, len(resp.Volumes))

    for _, v := range resp.Volumes {
        result = append(result, Volume{
            Name:       v.Name,
            Driver:     v.Driver,
            Mountpoint: v.Mountpoint,
            Labels:     v.Labels,
            Scope:      v.Scope,
        })
    }

    return result, nil
}


// ------------------------------------------ CREATE VOLUME 
func (c *Client) CreateVolume(
    ctx context.Context,
    req CreateVolumeRequest,
) error {
		driver := req.Driver
		if driver == ""{
			driver = "local"
		}
    _, err := c.cli.VolumeCreate(
        ctx,
        volume.CreateOptions{
            Name:       req.Name,
            Driver:     driver,
            Labels:     req.Labels,
            DriverOpts: req.Options,
        },
    )

    return err
}


// ---------------------------------------------- REMOVE VOLUME
func (c *Client) RemoveVolume(
    ctx context.Context,
    name string,
    force bool,
) error {
    return c.cli.VolumeRemove(
        ctx,
        name,
        force,
    )
}
