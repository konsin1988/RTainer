package grpcserver

import (
    "context"
		"strings"

		"konsin1988/agent/docker"
    pb "konsin1988/agent/proto"
)


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


// ------------------------------------------------------------
// REMOVE IMAGE 
// -------------------------------------------------------------
func (s *Server) RemoveImage(
    ctx context.Context,
    req *pb.RemoveImageRequest,
) (*pb.RemoveImageResponse, error) {

    err := s.imageSvc.RemoveImage(
        ctx,
        req,
    )
    if err != nil {
        return nil, err
    }

    return &pb.RemoveImageResponse{}, nil
}

// ------------------------------------------------------
// INSPECT IMAGE 
// ------------------------------------------------------
func (s *Server) InspectImage(
    ctx context.Context,
    req *pb.ImageRequest,
) (*pb.InspectImageResponse, error) {

    img, err := s.imageSvc.InspectImage(
        ctx,
        req.Id,
    )
    if err != nil {
        return nil, err
    }

    return &pb.InspectImageResponse {
        Id:            img.ID,
        RepoTags:      img.RepoTags,
        RepoDigests:   img.RepoDigests,
        Size:          img.Size,
        Created:       img.Created,
        Os:            img.OS,
        Architecture:  img.Architecture,
        Env:           img.Env,
        Cmd:           img.Cmd,
        Entrypoint:    img.Entrypoint,
        Labels:        img.Labels,
        ExposedPorts:  img.ExposedPorts,
    }, nil
}


// --------------------------------------------- PULL IMAGE
func (s *Server) PullImage(
    req *pb.PullImageRequest,
    stream pb.AgentService_PullImageServer,
) error {

    return s.imageSvc.PullImage(
        stream.Context(),
        req,
        func(p docker.PullProgress) error {

            return stream.Send(&pb.PullImageMessage{
                Status:  p.Status,
                Id:      p.ID,
                Current: p.Current,
                Total:   p.Total,
                Error:   p.Error,
            })
        },
    )
}
