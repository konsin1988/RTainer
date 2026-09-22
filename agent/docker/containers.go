package docker

import (
    "context"
		"strings"

		"github.com/docker/docker/api/types/container"
		"github.com/docker/go-connections/nat"
		
		"konsin1988/agent/proto"
)


type Container struct {
    ID     string
    Name   string
    Image  string
    Status string
}

// ---------------------------------------------------- LIST CONTAINER 
func (c *Client) ListContainers(
    ctx context.Context,
) ([]container.Summary, error) {

  return c.cli.ContainerList(
      ctx,
      container.ListOptions{
          All: true,
      },
  )
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


// ---------------------------------------------------- STOP CONTAINER 
func (c *Client) StopContainer(
    ctx context.Context,
    id string,
) error {
		t := 10 
    return c.cli.ContainerStop(
        ctx,
        id,
        container.StopOptions{
            Timeout: &t,
        },
    )
}

// ---------------------------------------------------- START CONTAINER 
func (c *Client) StartContainer(
    ctx context.Context,
    id string,
) error {
    return c.cli.ContainerStart(
        ctx,
        id,
        container.StartOptions{},
    )
}

// ------------------------------------------------------ PAUSE CONTAINER 
func (c *Client) PauseContainer(
    ctx context.Context,
    id string,
) error {
    return c.cli.ContainerPause(
        ctx,
        id,
    )
}

// ------------------------------------------------------ UNPAUSE CONTAINER 
func (c *Client) UnpauseContainer(
    ctx context.Context,
    id string,
) error {
    return c.cli.ContainerUnpause(
        ctx,
        id,
    )
}

// ----------------------------------------------------- RESTART CONTAINER 
func (c *Client) RestartContainer(
    ctx context.Context,
    id string,
) error {
    return c.cli.ContainerRestart(
        ctx,
        id,
        container.StopOptions{},
    )
}

// ----------------------------------------------------- REMOVE CONTAINER 
func (c *Client) RemoveContainer(
	ctx context.Context,
	id string,
	force bool,
	removeVolumes bool,
) error {

	return c.cli.ContainerRemove(
		ctx,
		id,
		container.RemoveOptions{
			Force:         force,
			RemoveVolumes: removeVolumes,
		},
	)
}

// ----------------------------------------------------- KILL CONTAINER 
func (c *Client) KillContainer(
	ctx context.Context,
	id string,
	signal string,
) error {

	return c.cli.ContainerKill(
		ctx,
		id,
		signal,
	)
}

// ------------------------------------------------------ UPDATE CONTAINER 
func (c *Client) UpdateContainer(
    ctx context.Context,
		containerID string,
    update container.UpdateConfig,
) error {
		_, err := c.cli.ContainerUpdate(ctx, containerID, update)
		if err != nil {
			return err
		}
		return nil
}

// ------------------------------------------------------ RUN CONTAINER 
func (c *Client) RunContainer(
    ctx context.Context,
    req *proto.RunContainerRequest,
) error {

		exposedPorts := nat.PortSet{}
		portBindings := nat.PortMap{}
		
		for _, p := range req.Ports {
				containerPort := p.ContainerPort
				if !strings.Contains(containerPort, "/") {
				    containerPort += "/tcp"
				}
		    port := nat.Port(p.ContainerPort) // e.g. "80/tcp"
		
		    exposedPorts[port] = struct{}{}
		
		    portBindings[port] = []nat.PortBinding{
		        {
		            HostIP:   "",
		            HostPort: p.HostPort,
		        },
		    }
		}

		binds := make([]string, 0, len(req.Volumes))
		
		for _, v := range req.Volumes {
		    binds = append(binds, v.Source+":"+v.Target)
		}

    cfg := &container.Config{
        Image: req.ImageId,
        Cmd:   req.Command,
        Env:   req.Env,
        Tty:   req.Tty,
				ExposedPorts: exposedPorts,
    }

    hostCfg := &container.HostConfig{
			PortBindings: portBindings,
			Binds:        binds,
		}

    resp, err := c.cli.ContainerCreate(
        ctx,
        cfg,
        hostCfg,
        nil,
        nil,
        req.Name,
    )
    if err != nil {
        return err
    }

    err = c.cli.ContainerStart(
        ctx,
        resp.ID,
        container.StartOptions{},
    )
    if err != nil {
        return err
    }

    return  nil
}


// ------------------------------------------------------ RENAME CONTAINER 
func (c *Client) RenameContainer(
    ctx context.Context,
    id string,
		name string,
) error {
    return c.cli.ContainerRename (
        ctx,
        id,
				name,
    )
}
