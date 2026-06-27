package service

import (
    "context"

    "konsin1988/agent/docker"
)

type ContainerService struct {
    docker *docker.Client
}

func NewContainerService(d *docker.Client) *ContainerService {
    return &ContainerService{docker: d}
}

type Container struct {
    ID     string
    Name   string
    Image  string
    Status string
}

func (s *ContainerService) ListContainers(ctx context.Context) ([]Container, error) {
    ctrs, err := s.docker.ListContainers(ctx)
    if err != nil {
        return nil, err
    }

    // Example business rule (this is WHY this layer exists)
    filtered := make([]Container, 0, len(ctrs))

    for _, c := range ctrs {
        // example rule: ignore dead containers
        if c.Status == "Dead" {
            continue
        }

        filtered = append(filtered, Container(c))
    }

    return filtered, nil
}


func (s *ContainerService) StopContainer(
    ctx context.Context,
    id string,
    timeout int,
) error {

    var t *int
    if timeout > 0 {
        t = &timeout
    }

    return s.docker.StopContainer(ctx, id, t)
}
