package service

import (
	"context"
	"strings"

	"konsin1988/agent/docker"
	"konsin1988/agent/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageService struct {
	docker *docker.Client
}

func NewImageService(d *docker.Client) *ImageService {
	return &ImageService{docker: d}
}

type Image struct {
	ID       string
	RepoTags []string
	Size     int64
	Created  int64
}

// -------------------------------------------- LIST IMAGES
func (s *ImageService) ListImages(
	ctx context.Context,
) ([]Image, error) {
	images, err := s.docker.ListImages(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]Image, 0, len(images))

	for _, c := range images {
		filtered = append(filtered, Image(c))
	}

	return filtered, nil
}

// --------------------------------------------- REMOVE IMAGE
func (s *ImageService) RemoveImage(
	ctx context.Context,
	req *proto.RemoveImageRequest,
) error {

	return s.docker.RemoveImage(
		ctx,
		req.Id,
		req.Force,
	)
}

// --------------------------------------------- INSPECT IMAGE
func (s *ImageService) InspectImage(
	ctx context.Context,
	id string,
) (docker.ImageInfo, error) {

	return s.docker.InspectImage(ctx, id)
}

// ---------------------------------------------- PULL IMAGE
func (s *ImageService) PullImage(
	ctx context.Context,
	req *proto.PullImageRequest,
	onProgress func(docker.PullProgress) error,
) error {

	return s.docker.PullImage(
		ctx,
		docker.PullImageRequest{
			Reference: req.Reference,
		},
		onProgress,
	)
}

// ------------------------------------------------ TAG IMAGE
func (s *ImageService) TagImage(
	ctx context.Context,
	req *proto.TagImageRequest,
) (*proto.TagImageResponse, error) {

	if req.GetImageId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"id is required",
		)
	}
	if req.GetRepository() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"repository is required",
		)
	}
	if req.GetTag() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"tag is required",
		)
	}

	err := s.docker.TagImage(ctx, req.ImageId, req.Repository, req.Tag)
	if err != nil {
		return nil, err
	}
	resp := proto.TagImageResponse{
		ImageId:    req.ImageId,
		Repository: req.Repository,
		Tag:        req.Tag,
	}

	return &resp, nil
}

// -------------------------------------------------------------------- BUILD IMAGE
type BuildImageEventFunc func(*proto.BuildImageEvent) error

func (s *ImageService) BuildImage(
	ctx context.Context,
	req *proto.BuildImageRequest,
	send BuildImageEventFunc,
) error {
	imageName := strings.TrimSpace(req.GetImageName())
	if imageName == "" {
		return status.Error(
			codes.InvalidArgument,
			"image_name is required",
		)
	}

	dockerfile := strings.TrimSpace(req.GetDockerfile())
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}

	if len(req.GetContextTar()) == 0 {
		return status.Error(
			codes.InvalidArgument,
			"context_tar is required",
		)
	}

	tags := req.GetTags()
	if len(tags) == 0 {
		tags = []string{"latest"}
	}

	buildTags := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return status.Error(
				codes.InvalidArgument,
				"image tags cannot contain empty values",
			)
		}

		buildTags = append(buildTags, imageName+":"+tag)
	}

	buildArgs := ConvertBuildArgs(req.GetBuildArgs())

	_, err := s.docker.BuildImage(
		ctx,
		req.GetContextTar(),
		dockerfile,
		buildTags,
		req.GetNoCache(),
		req.GetPull(),
		req.GetBuildArgs(),
		send,
	)
}

func ConvertBuildArgs(
	args map[string]string,
) map[string]*string {
	result := make(map[string]*string, len(args))

	for key, value := range args {
		v := value
		result[key] = &v
	}

	return result
}
