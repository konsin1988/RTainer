package docker

import (
    "context"
		"io"
)

type PortBinding struct {
    ContainerPort string
		HostIP				string
    HostPort      string
}

type VolumeBinding struct {
    Source string
    Target string
}

type ContainerInspect struct {
    ID     string
    Name   string
    Image  string
    Status string

    Ports  []PortBinding
    Mounts []VolumeBinding
    Env    []string

    Health *HealthStatus
}

type HealthStatus struct {
    Status        string
    FailingStreak int
    Logs          []string
}



// ----------------------------- INSPECT CONTAINER ----------------
func (c *Client) InspectContainer(
    ctx context.Context,
    id string,
) (ContainerInspect, error) {

    result, err := c.cli.ContainerInspect(
        ctx,
        id,
    )

    if err != nil {
        return ContainerInspect{}, err
    }


    info := ContainerInspect{
        ID:     result.ID,
        Name:   result.Name,
        Image:  result.Config.Image,
        Status: result.State.Status,
        Env:    result.Config.Env,
    }

    // Ports
    for containerPort, bindings := range result.NetworkSettings.Ports {
        for _, binding := range bindings {
            info.Ports = append(info.Ports, PortBinding{
                ContainerPort: string(containerPort),
								HostIP: binding.HostIP,
                HostPort:      binding.HostPort,
            })
        }
    }

    // Mounts
    for _, mount := range result.Mounts {

        info.Mounts = append(info.Mounts, VolumeBinding{
            Source: mount.Source,
            Target: mount.Destination,
        })
    }


    // Health
    if result.State.Health != nil {

        health := &HealthStatus{
            Status:        result.State.Health.Status,
            FailingStreak: result.State.Health.FailingStreak,
        }

        for _, log := range result.State.Health.Log {
            health.Logs = append(
                health.Logs,
                log.Output,
            )
        }

        info.Health = health
    }
    return info, nil
}


// ----------------------- CONTAINER STATS -------------------
func (c *Client) ContainerStats(
    ctx context.Context,
    id string,
) (io.ReadCloser, error) {

    resp, err := c.cli.ContainerStats(ctx, id, true)
    if err != nil {
        return nil, err
    }

    return resp.Body, nil
}
