package grpcserver

import (
		"bufio"
		"io"
		"log"
		"sync"
		"encoding/json"
		"errors"
    "context"

		"github.com/docker/docker/pkg/stdcopy"
		"github.com/docker/docker/api/types/container"
		"google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "konsin1988/agent/proto"
)


// -------------------------------------------------------------------- LIST CONTAINERS 
func (s *Server) ListContainers(
    ctx context.Context,
    req *pb.ListContainersRequest,
) (*pb.ListContainersResponse, error) {

    ctrs, err := s.containerSvc.ListContainers(ctx)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListContainersResponse{}

    for _, c := range ctrs {
        resp.Containers = append(resp.Containers, &pb.Container{
            Id:     c.ID,
            Name:   c.Name,
            Image:  c.Image,
            Status: c.Status,
        })
    }

    return resp, nil
}

// --------------------------------------------------------------------- INSPECT CONTAINER 
func (s *Server) InspectContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil 
}

// ----------------------------------------------------------------------- CREATE CONTAINER
func (s *Server) CreateContainer(
    ctx context.Context,
    req *pb.RunContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.RunContainer(ctx, req)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}

// ----------------------------------------------------------------------- START CONTAINER
func (s *Server) StartContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

    err := s.containerSvc.StartContainer(ctx, req.Id )
    if err != nil {
        return nil, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}

// ---------------------------------------------------------------------- STOP CONTAINER
func (s *Server) StopContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

    err := s.containerSvc.StopContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}



// ---------------------------------------------------------------------- RESTART CONTAINER 
func (s *Server) RestartContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

    err := s.containerSvc.RestartContainer(
        ctx,
        req.Id,
    )
    if err != nil {
        return nil, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}

// ------------------------------------------------------------------- PAUSE CONTAINER
func (s *Server) PauseContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

    err := s.containerSvc.PauseContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}


// -------------------------------------------------------------------- UNPAUSE CONTAINER
func (s *Server) UnpauseContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

    err := s.containerSvc.UnpauseContainer(ctx, req.Id )
    if err != nil {
        return nil, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}

// ----------------------------------------------------------------------- REMOVE CONTAINER
func (s *Server) RemoveContainer(
    ctx context.Context,
    req *pb.RemoveContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "id is required",
        )
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

    err = s.containerSvc.RemoveContainer(ctx, req.Id, req.Force, req.RemoveVolumes)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

		return &pb.ContainerResponse{Container: resp}, nil
}


// ----------------------------------------------------------------- KILL CONTAINER
func (s *Server) KillContainer(
    ctx context.Context,
    req *pb.KillContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "container_id is required",
        )
    }

    signal := req.GetSignal()

    if signal == "" {
        signal = "SIGKILL"
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

    err = s.containerSvc.KillContainer(ctx, req.Id, signal)
    if err != nil {
        return nil, err
    }

		return &pb.ContainerResponse{Container: resp}, nil
}

// ----------------------------------------------------------------- UPDATE CONTAINER
func (s *Server) UpdateContainer(
    ctx context.Context,
    req *pb.UpdateContainerRequest,
) (*pb.ContainerResponse, error) {

		if req.GetId() == "" {
        return nil, status.Error(
            codes.InvalidArgument,
            "container_id is required",
        )
    }

    err := s.containerSvc.UpdateContainer(ctx, req)
    if err != nil {
        return nil, err
    }

		resp, err := s.containerSvc.InspectContainer(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		return &pb.ContainerResponse{Container: resp}, nil
}

// ------------------------------------------------------------------------- EXEC 
func (s *Server) ExecContainer (
    req *pb.ExecuteCommandRequest,
    stream pb.ContainerService_ExecContainerServer,
) error {

    result, err := s.containerSvc.ExecuteCommand(
        stream.Context(),
        req,
    )
    if err != nil {
        return err
    }
    defer result.Close()

    return s.streamLogMessages(
        stream.Context(),
        result.Reader,
        req.Tty,
        stream.Send,
    )
}

// --------------------------------------------------------------- LOGS CONTAINER 
func (s *Server) LogsContainer(
    req *pb.ViewLogsRequest,
    stream pb.ContainerService_LogsContainerServer,
) error {

    reader, err := s.containerSvc.ViewLogs(
        stream.Context(),
        req,
    )
    if err != nil {
        return err
    }
    defer reader.Close()

    return s.streamLogMessages(
        stream.Context(),
        reader,
        false, // Docker logs are always multiplexed
        stream.Send,
    )
}


// ----------------------------------------------------------------- STATS CONTAINER
func (s *Server) StatsContainer(
    req *pb.ContainerRequest,
    stream pb.ContainerService_StatsContainerServer,
) error {

    reader, err := s.containerSvc.ContainerStats(
        stream.Context(),
        req,
    )
    if err != nil {
        return err
    }
    defer reader.Close()

    decoder := json.NewDecoder(reader)

    var previous *container.StatsResponse

    for {

        var stats container.StatsResponse

        if err := decoder.Decode(&stats); err != nil {
            if errors.Is(err, io.EOF) {
                return nil
            }
            return err
        }

        var cpuPercent float64

        if previous != nil {

            cpuDelta :=
                float64(stats.CPUStats.CPUUsage.TotalUsage -
                    previous.CPUStats.CPUUsage.TotalUsage)

            systemDelta :=
                float64(stats.CPUStats.SystemUsage -
                    previous.CPUStats.SystemUsage)

            if cpuDelta > 0 && systemDelta > 0 {

                onlineCPUs := stats.CPUStats.OnlineCPUs
                if onlineCPUs == 0 {
                    onlineCPUs = uint32(len(stats.CPUStats.CPUUsage.PercpuUsage))
                }

                cpuPercent =
                    (cpuDelta / systemDelta) *
                        float64(onlineCPUs) *
                        100.0
            }
        }

        previous = &stats

        var rx uint64
        var tx uint64

        for _, network := range stats.Networks {
            rx += network.RxBytes
            tx += network.TxBytes
        }

        var blockRead uint64
        var blockWrite uint64

        for _, io := range stats.BlkioStats.IoServiceBytesRecursive {

            switch io.Op {

            case "Read":
                blockRead += io.Value

            case "Write":
                blockWrite += io.Value
            }
        }

        err = stream.Send(&pb.ContainerStatsResponse{
            CpuPercent: cpuPercent,

            MemoryUsage: stats.MemoryStats.Usage,
            MemoryLimit: stats.MemoryStats.Limit,

            NetworkRx: rx,
            NetworkTx: tx,

            BlockRead: blockRead,
            BlockWrite: blockWrite,

            Pids: uint32(stats.PidsStats.Current),
        })

        if err != nil {
            return err
        }
    }
}


// --------------------------------------
// VIEW LOGS AND EXECUTE COMMAND HELPER
// ---------------------------------------

type logEntry struct {
    line   string
    stream pb.LogStream
}

func (s *Server) streamLogMessages(
    ctx context.Context,
    reader io.Reader,
    tty bool,
    send func(*pb.LogMessage) error,
) error {

    // TTY: stdout/stderr are merged.
    if tty {
        scanner := bufio.NewScanner(reader)

        for scanner.Scan() {
            if err := send(&pb.LogMessage{
                Line:   scanner.Text(),
                Stream: pb.LogStream_STDOUT,
            }); err != nil {
                return err
            }
        }

        return scanner.Err()
    }

    stdoutReader, stdoutWriter := io.Pipe()
    stderrReader, stderrWriter := io.Pipe()

    logs := make(chan logEntry)

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer stdoutWriter.Close()
        defer stderrWriter.Close()

        _, err := stdcopy.StdCopy(stdoutWriter, stderrWriter, reader)
        if err != nil {
            stdoutWriter.CloseWithError(err)
            stderrWriter.CloseWithError(err)
        }
    }()

    scan := func(
        r io.Reader,
        streamType pb.LogStream,
    ) {
        defer wg.Done()

        scanner := bufio.NewScanner(r)

        for scanner.Scan() {
            select {
            case logs <- logEntry{
                line:   scanner.Text(),
                stream: streamType,
            }:
            case <-ctx.Done():
                return
            }
        }

        if err := scanner.Err(); err != nil && ctx.Err() == nil {
            log.Printf("scanner error: %v", err)
        }
    }

    go scan(stdoutReader, pb.LogStream_STDOUT)
    go scan(stderrReader, pb.LogStream_STDERR)

    go func() {
        wg.Wait()
        close(logs)
    }()

    for entry := range logs {
        if err := send(&pb.LogMessage{
            Line:   entry.line,
            Stream: entry.stream,
        }); err != nil {
            return err
        }
    }

    return nil
}

