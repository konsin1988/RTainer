package grpcserver

import (
		"bufio"
    "context"
		"io"
		"log"
		"sync"
		"encoding/json"
		"errors"

		"github.com/docker/docker/pkg/stdcopy"
		"github.com/docker/docker/api/types/container"

    pb "konsin1988/agent/proto"
)

// --------------------------------------
// VIEW LOGS 
// ---------------------------------------

type logEntry struct {
    line   string
    stream pb.LogStream
}

func (s *Server) ViewLogs(
    req *pb.ViewLogsRequest,
    stream pb.AgentService_ViewLogsServer,
) error {

    reader, err := s.containerSvc.ViewLogs(stream.Context(), req)
    if err != nil {
        return err
    }
    defer reader.Close()

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

		go func() {
				defer wg.Done()

		    scanner := bufio.NewScanner(stdoutReader)
		
		    for scanner.Scan() {
						select {
						    case logs <- logEntry{
						        line:   scanner.Text(),
						        stream: pb.LogStream_STDOUT,
						    }:
						    case <-stream.Context().Done():
						        return
						    }
		    }

				if err := scanner.Err(); err != nil && stream.Context().Err() == nil {
				    log.Printf("stdout scanner error: %v", err)
				}
		}()

		go func() {
				defer wg.Done()

				scanner := bufio.NewScanner(stderrReader)

				for scanner.Scan() {
						select {
						case logs <- logEntry{
						    line: scanner.Text(),
						    stream: pb.LogStream_STDERR,
						}:
						case <-stream.Context().Done():
						    return
						}
    		}

				if err := scanner.Err(); err != nil && stream.Context().Err() == nil {
				    log.Printf("stderr scanner error: %v", err)
				}
		}()


		go func() {
		    wg.Wait()
		    close(logs)
		}()

		for entry := range logs {
		    err := stream.Send(&pb.LogMessage{
		        Line:   entry.line,
		        Stream: entry.stream,
		    })
		    if err != nil {
		        return err
		    }
		}

		return nil
}


// ----------------------------------------
// INSPECT CONTAINER 
// ---------------------------------------------
func (s *Server) InspectContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.InspectContainerResponse, error) {


    result, err := s.containerSvc.InspectContainer(
        ctx,
        req,
    )

    if err != nil {
        return nil, err
    }

    resp := &pb.InspectContainerResponse{
        Id:     result.ID,
        Name:   result.Name,
        Image:  result.Image,
        Status: result.Status,
        Env:    result.Env,
    }

    for _, p := range result.Ports {

        resp.Ports = append(
            resp.Ports,
            &pb.PortBinding{
                ContainerPort: p.ContainerPort,
								HostIp: 			 p.HostIP,
                HostPort:      p.HostPort,
            },
        )
    }

    for _, v := range result.Mounts {

        resp.Mounts = append(
            resp.Mounts,
            &pb.VolumeBinding{
                Source: v.Source,
                Target: v.Target,
            },
        )
    }
    if result.Health != nil {

        resp.Health = &pb.HealthStatus{
            Status:        result.Health.Status,
            FailingStreak: int32(result.Health.FailingStreak),
            Logs:          result.Health.Logs,
        }
    }


    return resp, nil
}


// -----------------------------------------------
// CONTAINER STATS 
// --------------------------------------------------
func (s *Server) ContainerStats(
    req *pb.ContainerRequest,
    stream pb.AgentService_ContainerStatsServer,
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
