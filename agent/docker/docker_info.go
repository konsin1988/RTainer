package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
)

type DockerInfo struct {
    ServerVersion string
    APIVersion    string

    OperatingSystem string
    Architecture    string
    KernelVersion   string

    Containers        int32
    ContainersRunning int32
    ContainersPaused  int32
    ContainersStopped int32

    Images  int32
    Volumes int32

    Driver string

    MemTotal int64
    NCPU     int32

    DockerRootDir string

    Name string
    ID   string
}


func (c *Client) DockerInfo(
    ctx context.Context,
) (DockerInfo, error) {

    info, err := c.cli.Info(ctx)
    if err != nil {
        return DockerInfo{}, err
    }

		volumes, err := c.cli.VolumeList(ctx, volume.ListOptions{})
		if err != nil {
		    return DockerInfo{}, err
		}

    version, err := c.cli.ServerVersion(ctx)
    if err != nil {
        return DockerInfo{}, err
    }

    return DockerInfo{
        ServerVersion: version.Version,
        APIVersion:    version.APIVersion,

        OperatingSystem: info.OperatingSystem,
        Architecture:    info.Architecture,
        KernelVersion:   info.KernelVersion,

        Containers:        int32(info.Containers),
        ContainersRunning: int32(info.ContainersRunning),
        ContainersPaused:  int32(info.ContainersPaused),
        ContainersStopped: int32(info.ContainersStopped),

        Images: int32(info.Images),
        Volumes: int32(len(volumes.Volumes)),

        Driver: info.Driver,

        MemTotal: info.MemTotal,
        NCPU:     int32(info.NCPU),

        DockerRootDir: info.DockerRootDir,

        Name: info.Name,
        ID:   info.ID,
    }, nil
}
