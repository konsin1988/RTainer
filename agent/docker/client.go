package docker

import (
		dockerclient "github.com/docker/docker/client"
)

type Client struct {
    cli *dockerclient.Client
}

func New() (*Client, error) {

    cli, err := dockerclient.NewClientWithOpts(
        dockerclient.FromEnv,
				dockerclient.WithAPIVersionNegotiation(),
    )

    if err != nil {
        return nil, err
    }

    return &Client{
        cli: cli,
    }, nil
}

