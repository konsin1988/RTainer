package docker

import (
    "context"
		"io"
		"strconv"

		"github.com/docker/docker/api/types/container"
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


// -------------------------------- VIEW LOGS ----------------
type LogsRequest struct {
    ContainerID string
    Follow      bool
    Tail        int
    Timestamps  bool
}

func (c *Client) ContainerLogs(
    ctx context.Context,
		req LogsRequest,
) (io.ReadCloser, error) {

		tail := "all"
		
		if req.Tail > 0 {
		    tail = strconv.Itoa(req.Tail)
		}

    return c.cli.ContainerLogs(ctx, req.ContainerID, 
		container.LogsOptions{
            ShowStdout: true,
            ShowStderr: true,
            Follow:     req.Follow,
            Tail:       tail,
            Timestamps: req.Timestamps,
        })
}

// ------------------------------ EXECUTE COMMAND ------------------
type ExecRequest struct {
	ContainerID string
	Command     []string
	Tty         bool
}

type ExecResult struct {
    Reader 	io.Reader
		Close 	func() 
}

func (c *Client) ExecuteCommand(
	ctx context.Context,
	req ExecRequest,
) (ExecResult, error) {

	// Create exec instance
	execResp, err := c.cli.ContainerExecCreate(
		ctx,
		req.ContainerID,
		container.ExecOptions{
			Cmd:          req.Command,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          req.Tty,
		},
	)

	if err != nil {
      return ExecResult{}, err
  }

	resp, err := c.cli.ContainerExecAttach(
		ctx,
		execResp.ID,
		container.ExecAttachOptions{},
	)

	return ExecResult{
	    Reader: resp.Reader,
	    Close: resp.Close,
	}, nil
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
