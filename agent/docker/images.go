package docker

import (
    "context"

		"github.com/docker/docker/api/types/image"
)

type Image struct {
    ID        string
    RepoTags  []string
    Size      int64
    Created   int64
}


func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
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
