package docker

import (
    "context"
		"encoding/json"
		"io"

		"github.com/docker/docker/api/types/image"
		"github.com/docker/docker/pkg/jsonmessage"
)

type Image struct {
    ID        string
    RepoTags  []string
    Size      int64
    Created   int64
}

type ImageInfo struct {
    ID            string
    RepoTags      []string
    RepoDigests   []string

    Size          int64
    Created       int64

    OS            string
    Architecture  string

    Env           []string
    Cmd           []string
    Entrypoint    []string

    Labels        map[string]string

    ExposedPorts  []string
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
