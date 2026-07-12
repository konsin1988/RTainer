package grpcserver

import (
		"bufio"
    "context"
		"io"
		"log"
		"strings"
		"sync"

    pb "konsin1988/agent/proto"
    "konsin1988/agent/service"

		"github.com/docker/docker/pkg/stdcopy"
    "google.golang.org/grpc"
)

type Server struct {
    pb.UnimplementedAgentServiceServer
    containerSvc *service.ContainerService
		imageSvc *service.ImageService
}

func New(container *service.ContainerService, image *service.ImageService) *Server {
		return &Server{containerSvc: container, imageSvc: image}
}

func (s *Server) Register(grpcSrv *grpc.Server) {
    pb.RegisterAgentServiceServer(grpcSrv, s)
}


// ---------------------------------------
// LIST CONTAINERS 
// --------------------------------------
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

// -----------------------------------
// LIST IMAGES
// ------------------------------------
func (s *Server) ListImages(
		ctx context.Context,
		req *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
    images, err := s.imageSvc.ListImages(ctx)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListImagesResponse{}

    for _, c := range images {
        resp.Images = append(resp.Images, &pb.Image{
            Id:     	c.ID,
						RepoTags: strings.Join(c.RepoTags, ","), 
						Size:			c.Size,     
						Created: 	c.Created,   
        })
    }

    return resp, nil
}

// ------------------------------------
// STOP CONTAINER
// ------------------------------------
func (s *Server) StopContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.StopContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}


// ------------------------------------
// START CONTAINER
// ------------------------------------
func (s *Server) StartContainer(
    ctx context.Context,
    req *pb.ContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.StartContainer(ctx, req.Id )
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}


// ------------------------------------
// REMOVE CONTAINER
// ------------------------------------
func (s *Server) RemoveContainer(
    ctx context.Context,
    req *pb.RemoveContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.RemoveContainer(ctx, req.Id, req.Force, req.RemoveVolumes)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}


// ------------------------------------
// RUN CONTAINER
// ------------------------------------
func (s *Server) RunContainer(
    ctx context.Context,
    req *pb.RunContainerRequest,
) (*pb.ContainerResponse, error) {

    err := s.containerSvc.RunContainer(ctx, req)
    if err != nil {
        return &pb.ContainerResponse{}, err
    }

    return &pb.ContainerResponse{}, nil
}

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
