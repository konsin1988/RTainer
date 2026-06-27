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
    timeout *int,
) error {

		t := 10 // default timeout

    if timeout != nil {
				t = *timeout
    }

    return c.cli.ContainerStop(
        ctx,
        id,
        container.StopOptions{
            Timeout: &t,
        },
    )
}
