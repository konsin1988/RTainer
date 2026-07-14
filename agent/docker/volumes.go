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
