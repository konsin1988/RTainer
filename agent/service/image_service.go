package service

import (
    "context"

    "konsin1988/agent/docker"
)

type ImageService struct {
    docker *docker.Client
}

func NewImageService(d *docker.Client) *ImageService {
    return &ImageService{docker: d}
}

type Image struct {
    ID        string
    RepoTags  []string
    Size      int64
    Created   int64
}

func (s *ImageService) ListImages(ctx context.Context) ([]Image, error) {
    images, err := s.docker.ListImages(ctx)
    if err != nil {
        return nil, err
    }

    filtered := make([]Image, 0, len(images))

    for _, c := range images {
        // example rule: ignore dead containers
        //if c.Status == "Dead" {
        //    continue
        //}

        filtered = append(filtered, Image(c))
    }

    return filtered, nil
}
