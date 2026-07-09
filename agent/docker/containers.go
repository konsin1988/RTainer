package docker

import (
    "context"

		"github.com/docker/docker/api/types/container"
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
