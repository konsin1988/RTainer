package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
)

type Image struct {
	ID       string
	RepoTags []string
	Size     int64
	Created  int64
}

type ImageInfo struct {
	ID          string
	RepoTags    []string
	RepoDigests []string

	Size    int64
	Created int64

	OS           string
	Architecture string

	Env        []string
	Cmd        []string
	Entrypoint []string

	Labels map[string]string

	ExposedPorts []string
}

type PullImageRequest struct {
	Reference string
}

type PullProgress struct {
	Status  string
	ID      string
	Current int64
	Total   int64
	Error   string
}

// ---------------------------------------------------- LIST IMAGE
func (c *Client) ListImages(
	ctx context.Context,
) ([]Image, error) {
	imgs, err := c.cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]Image, 0, len(imgs))

	for _, img := range imgs {
		result = append(result, Image{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
		})
	}

	return result, nil
}

// ------------------------------------------------- REMOVE IMAGE
func (c *Client) RemoveImage(
	ctx context.Context,
	id string,
	force bool,
) error {

	_, err := c.cli.ImageRemove(
		ctx,
		id,
		image.RemoveOptions{
			Force: force,
		},
	)

	return err
}

// ------------------------------------------------- INSPECT IMAGE
func (c *Client) InspectImage(
	ctx context.Context,
	id string,
) (ImageInfo, error) {

	img, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return ImageInfo{}, err
	}

	info := ImageInfo{
		ID:           img.ID,
		RepoTags:     img.RepoTags,
		RepoDigests:  img.RepoDigests,
		Size:         img.Size,
		OS:           img.Os,
		Architecture: img.Architecture,
		Labels:       img.Config.Labels,
	}

	if img.Config != nil {
		info.Env = img.Config.Env
		info.Cmd = img.Config.Cmd
		info.Entrypoint = img.Config.Entrypoint

		for p := range img.Config.ExposedPorts {
			info.ExposedPorts = append(
				info.ExposedPorts,
				string(p),
			)
		}
	}

	return info, nil
}

// ----------------------------------------------------PULL IMAGE
func (c *Client) PullImage(
	ctx context.Context,
	req PullImageRequest,
	onProgress func(PullProgress) error,
) error {

	reader, err := c.cli.ImagePull(
		ctx,
		req.Reference,
		image.PullOptions{},
	)
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)

	for {
		var msg jsonmessage.JSONMessage

		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		progress := PullProgress{
			Status: msg.Status,
			ID:     msg.ID,
		}

		if err := msg.Error; err != nil {
			progress.Error = err.Message
		}

		if msg.Progress != nil {
			progress.Current = msg.Progress.Current
			progress.Total = msg.Progress.Total
		}

		if err := onProgress(progress); err != nil {
			return err
		}
	}
}

// ------------------------------------------------------------------ TAG IMAGE
func (c *Client) TagImage(ctx context.Context, id, repository, tag string) error {
	target := repository + ":" + tag
	return c.cli.ImageTag(ctx, id, target)
}

// ----------------------------------------------------------------- BUILD IMAGE
type BuildImageEvent struct {
	Stream  string
	Error   string
	ImageID string
}

type BuildImageEventFunc func(BuildImageEvent) error

type dockerBuildMessage struct {
	Stream      string `json:"stream"`
	Error       string `json:"error"`
	ErrorDetail struct {
		Message string `json:"message"`
	} `json:"errorDetail"`

	Aux struct {
		ID string `json:"ID"`
	} `json:"aux"`
}

func (c *Client) BuildImage(
	ctx context.Context,
	contextTar []byte,
	dockerfile string,
	tags []string,
	noCache bool,
	pull bool,
	buildArgs map[string]string,
	send BuildImageEventFunc,
) (string, error) {

	args := make(map[string]*string, len(buildArgs))
	for key, value := range buildArgs {
		v := value
		args[key] = &v
	}

	options := build.ImageBuildOptions{
		Dockerfile: dockerfile,
		Tags:       tags,
		Remove:     true,
		NoCache:    noCache,
		PullParent: pull,
		BuildArgs:  args,
	}

	response, err := c.cli.ImageBuild(
		ctx,
		bytes.NewReader(contextTar),
		options,
	)
	if err != nil {
		return "", fmt.Errorf("start image build: %w", err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var imageID string

	for {
		var msg dockerBuildMessage
		err := decoder.Decode(&msg)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return imageID, fmt.Errorf("decode Docker build output: %w", err)
		}

		if msg.Aux.ID != "" {
			imageID = msg.Aux.ID
		}

		buildErr := msg.Error
		if buildErr == "" {
			buildErr = msg.ErrorDetail.Message
		}
		if buildErr != "" {
			if send != nil {
				if err := send(BuildImageEvent{Error: buildErr}); err != nil {
					return imageID, fmt.Errorf("send build error: %w", err)
				}
			}
			return imageID, fmt.Errorf("Docker image build failed: %s", buildErr)
		}
	}

	return imageID, nil
}
