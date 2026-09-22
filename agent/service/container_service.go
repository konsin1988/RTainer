package service

import (
    "context"
		"io"
		"strings"

		"github.com/docker/docker/api/types/container"
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
) (*proto.ListContainersResponse, error) {

    ctrs, err := s.docker.ListContainers(ctx)
    if err != nil {
        return nil, err
    }

    resp := &proto.ListContainersResponse{}

    for _, c := range ctrs {
        resp.Containers = append(resp.Containers, &proto.Container{
            Id:     c.ID,
            Name:   containerName(c), 
            Image:  c.Image,
            Status: c.Status,
        })
    }

    return resp, nil
}

func containerName(c container.Summary) string {
    if len(c.Names) == 0 {
        return ""
    }

    return strings.TrimLeft(c.Names[0], "/")
}

// -------------------------------------------- STOP CONTAINER 
func (s *ContainerService) StopContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.StopContainer(ctx, id)
}


// --------------------------------------------------------------- START CONTAINER 
func (s *ContainerService) StartContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.StartContainer(ctx, id)
}

// ---------------------------------------------------------------- RESTART CONTAINER 
func (s *ContainerService) RestartContainer(
    ctx context.Context,
    id string,
) error {
    return s.docker.RestartContainer(ctx, id)
}

// ---------------------------------------------------------------- PAUSE CONTAINER 
func (s *ContainerService) PauseContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.PauseContainer(ctx, id)
}

// --------------------------------------------------------------- UNPAUSE CONTAINER 
func (s *ContainerService) UnpauseContainer(
    ctx context.Context,
    id string,
) error {

    return s.docker.UnpauseContainer(ctx, id)
}

// ----------------------------------------------------------------- REMOVE CONTAINER 
func (s *ContainerService) RemoveContainer(
		ctx context.Context,
		id string,
		force bool,
		removeVolumes bool,
) error {
    return s.docker.RemoveContainer(ctx, id, force, removeVolumes)
}

// ----------------------------------------------------------------------- KILL CONTAINER 
func (s *ContainerService) KillContainer(
		ctx context.Context,
		id string,
		signal string,
) error {
    return s.docker.KillContainer(ctx, id, signal)
}

// ----------------------------------------------------------------------- RUN CONTAINER 
func (s *ContainerService) RunContainer(
    ctx context.Context,
    req *proto.RunContainerRequest,
) error {
    return s.docker.RunContainer(ctx, req)
}

// ---------------------------------------------------------------- RENAME CONTAINER 
func (s *ContainerService) RenameContainer(
    ctx context.Context,
    id string,
		name string,
) error {

    return s.docker.RenameContainer(ctx, id, name)
}

// ---------------------------------------------------------------------	VIEW LOGS 
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

// ----------------------------------------------------------------------- UPDATE CONTAINER 
func (s *ContainerService) UpdateContainer(
    ctx context.Context,
    req *proto.UpdateContainerRequest,
) error {
		update := containerUpdateFromProto(req)
    return s.docker.UpdateContainer(ctx, req.Id, update)
}

func containerUpdateFromProto(
    req *proto.UpdateContainerRequest,
) container.UpdateConfig {

    update := container.UpdateConfig{}

    if req.CpuShares != nil {
        update.Resources.CPUShares = req.GetCpuShares()
    }

    if req.CpuPeriod != nil {
        update.Resources.CPUPeriod = req.GetCpuPeriod()
    }

    if req.CpuQuota != nil {
        update.Resources.CPUQuota = req.GetCpuQuota()
    }

    if req.Memory != nil {
        update.Resources.Memory = req.GetMemory()
    }

    if req.MemorySwap != nil {
        update.Resources.MemorySwap = req.GetMemorySwap()
    }

    if req.PidsLimit != nil {
        update.Resources.PidsLimit = req.PidsLimit
    }

    if req.BlkioWeight != nil {
        update.Resources.BlkioWeight = uint16(req.GetBlkioWeight())
    }

    if req.CpuCount != nil {
        update.Resources.CPUCount = req.GetCpuCount()
    }

    if req.CpuPercent != nil {
        update.Resources.CPUPercent = req.GetCpuPercent()
    }

    if req.RestartPolicy != nil {
        update.RestartPolicy = container.RestartPolicy{
            Name:              container.RestartPolicyMode(req.RestartPolicy.GetName()),
            MaximumRetryCount: int(req.RestartPolicy.GetMaximumRetryCount()),
        }
    }

    return update
}

// ------------------------------------------------------------------------ INSPECT CONTAINER 
func (s *ContainerService) InspectContainer(
    ctx context.Context,
    id string,
) (*proto.ContainerInfo, error) {

		result, err := s.docker.InspectContainer(
        ctx,
        id,
    )

		if err != nil {
			return nil, err
		}


    resp := &proto.ContainerInfo{
        Id:     result.ID,
        Name:   result.Name,
        Image:  result.Image,
        Status: result.Status,
        Env:    result.Env,
    }

    for _, p := range result.Ports {

        resp.Ports = append(
            resp.Ports,
            &proto.PortBinding{
                ContainerPort: p.ContainerPort,
								HostIp: 			 p.HostIP,
                HostPort:      p.HostPort,
            },
        )
    }

    for _, v := range result.Mounts {

        resp.Mounts = append(
            resp.Mounts,
            &proto.VolumeBinding{
                Source: v.Source,
                Target: v.Target,
            },
        )
    }
    if result.Health != nil {

        resp.Health = &proto.HealthStatus{
            Status:        result.Health.Status,
            FailingStreak: int32(result.Health.FailingStreak),
            Logs:          result.Health.Logs,
        }
    }


    return resp, nil
}

// -------------------------------------------------------------------- CONTAINER STATS 
func (s *ContainerService) ContainerStats(
    ctx context.Context,
    req *proto.ContainerRequest,
) (io.ReadCloser, error) {

    return s.docker.ContainerStats(
        ctx,
        req.Id,
    )
}


// -------------------------------------------------------------------- EXECUTE COMMAND 
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


// -------------------------------------------------------------------- DOCKER INFO
func (s *ContainerService) DockerInfo(
    ctx context.Context,
) (docker.DockerInfo, error) {
    return s.docker.DockerInfo(ctx)
}

// ----------------------------------------------------------------------- EVENTS
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
