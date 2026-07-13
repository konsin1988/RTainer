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

func (c *Client) ListContainers(
    ctx context.Context,
) ([]Container, error) {

  list, err := c.cli.ContainerList(
      ctx,
      container.ListOptions{
          All: true,
      },
  )

  if err != nil {
      return nil, err
  }

	containers := make([]Container, 0, len(list))

	for _, ctr := range list {
	    name := ""
	    if len(ctr.Names) > 0 {
	        name = ctr.Names[0]
	    }
	
	    containers = append(containers, Container{
	        ID:     ctr.ID,
	        Name:   name,
	        Image:  ctr.Image,
	        Status: ctr.Status,
	    })
	}
	
	return containers, nil
}

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

// -----------------------------------START CONTAINER ----------------
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

// ------------------------------- REMOVE CONTAINER ---------------
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


// -------------------------------- RUN CONTAINER ---------------------
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


