package service

import (
    "context"
		"io"

    "konsin1988/agent/docker"
		"konsin1988/agent/proto"
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

// -------------------------------------------------- LIST CONTAINERS
func (s *ContainerService) ListContainers(
	ctx context.Context,
) ([]Container, error) {
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


// -------------------------------------------- STOP CONTAINER 
func (s *ContainerService) StopContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.StopContainer(ctx, id)
}


// ------------------------------------------- START CONTAINER 
func (s *ContainerService) StartContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.StartContainer(ctx, id)
}

// -------------------------------------------- RESTART CONTAINER 
func (s *ContainerService) RestartContainer(
    ctx context.Context,
    id string,
) error {
    return s.docker.RestartContainer(ctx, id)
}

// -------------------------------------------- REMOVE CONTAINER 
func (s *ContainerService) RemoveContainer(
		ctx context.Context,
		id string,
		force bool,
		removeVolumes bool,
) error {
    return s.docker.RemoveContainer(ctx, id, force, removeVolumes)
}


// ------------------------------------------ RUN CONTAINER 
func (s *ContainerService) RunContainer(
    ctx context.Context,
    req *proto.RunContainerRequest,
) error {
    return s.docker.RunContainer(ctx, req)
}

// ------------------------------------------	VIEW LOGS 
func (s *ContainerService) ViewLogs(
	ctx context.Context,
	req *proto.ViewLogsRequest,
) (io.ReadCloser, error) {

    return s.docker.ContainerLogs(
      ctx,
			docker.LogsRequest{
        ContainerID: req.ContainerId,
        Follow:      req.Follow,
        Tail:        int(req.Tail),
        Timestamps:  req.Timestamps,
    	},
    )
}

// ------------------------------------------ INSPECT CONTAINER 
func (s *ContainerService) InspectContainer(
    ctx context.Context,
    req *proto.ContainerRequest,
) (docker.ContainerInspect, error) {

    return s.docker.InspectContainer(
        ctx,
        req.Id,
    )
}

// ----------------------------------------- CONTAINER STATS 
func (s *ContainerService) ContainerStats(
    ctx context.Context,
    req *proto.ContainerRequest,
) (io.ReadCloser, error) {

    return s.docker.ContainerStats(
        ctx,
        req.Id,
    )
}


// ---------------------------------------- EXECUTE COMMAND 
func (s *ContainerService) ExecuteCommand(
    ctx context.Context,
    req *proto.ExecuteCommandRequest,
) (docker.ExecResult, error) {

    return s.docker.ExecuteCommand(
        ctx,
        docker.ExecRequest{
            ContainerID: req.Id,
            Command:     req.Command,
            Tty:         req.Tty,
        },
    )
}


// ------------------------------------------- DOCKER INFO
func (s *ContainerService) DockerInfo(
    ctx context.Context,
) (docker.DockerInfo, error) {
    return s.docker.DockerInfo(ctx)
}

// ------------------------------------------------ EVENTS
func (s *ContainerService) Events(
    ctx context.Context,
    req *proto.EventsRequest,
) (<-chan docker.Event, <-chan error) {

    return s.docker.Events(
        ctx,
        docker.EventsRequest{
            Types: req.Types,
            Actions: req.Actions,
        },
    )
}
