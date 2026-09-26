package service

import (
	"context"
	"fmt"

	"konsin1988/agent/docker"
	"konsin1988/agent/proto"
		
	"github.com/docker/docker/api/types/container"
  _ "github.com/docker/docker/api/types/network"
  _ "github.com/docker/docker/api/types/volume"
)

type StackService struct {
  docker *docker.Client
	volume  *VolumeService
}

func NewStackService(d *docker.Client, volume *VolumeService) *StackService {
	return &StackService{docker: d, volume: volume}
}


// ------------------------------------------ LIST STACKS 
func (s *StackService) ListStacks(
    ctx context.Context,
) (*proto.ListStackResponse, error) {

	containers, err := s.docker.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]container.Summary)
	networkIDs:= make(map[string]map[string]bool)
	volumeIDs := make(map[string]struct{})
	networkGroups := make(map[string][]*proto.NetworkListItem)
    
	for _, c := range containers {
		stackName := c.Labels["com.docker.compose.project"]

    if stackName == "" {
        continue
    }

		if c.NetworkSettings == nil {
			continue
		}

		for _, network := range c.NetworkSettings.Networks {
				networkID := network.NetworkID
        if networkIDs[stackName] == nil {
            networkIDs[stackName] = make(map[string]bool)
        }

        networkIDs[stackName][networkID] = true
    }


    groups[stackName] = append(groups[stackName], c)
  }


	for stackName, names := range networkIDs{
	    for networkID := range names {
	
	        n, err := s.docker.InspectNetwork(ctx, networkID)
	        if err != nil {
	            return nil, err
	        }

					resp := &proto.NetworkListItem{
          	Id:     		n.ID,
          	Name:   		n.Name,
          	Driver: 		n.Driver,
          	Scope:  		n.Scope,
						Labels: 		n.Labels,
					}

					for _, c := range n.Containers{
						resp.Containers = append(
							resp.Containers,
							&proto.NetContainer{
								Id:						c.ID,
								Name:					c.Name,
                Ipv4Address: 	c.IPv4Address,
							})
					}
	
	        networkGroups[stackName] = append(
	            networkGroups[stackName],
							resp,
	        )
	    }
	}

	response := &proto.ListStackResponse{}

  for stackName, stackContainers := range groups {

    stack := &proto.Stack{
        Id:   stackName,
        Name: stackName,
    }

    for _, c := range stackContainers {
        stack.Containers = append(
            stack.Containers,
            &proto.Container{
                Id:     c.ID,
                Name:   containerName(c),
                Image:  c.Image,
                Status: c.Status,
            },
        )

				for _, mount := range c.Mounts {
				    fmt.Println(mount.Type)
				    fmt.Println(mount.Source)
				    fmt.Println(mount.Destination)
				}

				for _, mount := range c.Mounts {
    		    if mount.Type != "volume" {
    		        continue
    		    }

    		    if mount.Name == "" {
    		        continue
    		    }

    		    volumeIDs[mount.Name] = struct{}{}
    		}

				for volumeName := range volumeIDs {
  			  volume, err := s.docker.InspectVolume(ctx, volumeName)
  			  if err != nil {
  			      return nil, err
  			  }

					if volume.Labels["com.docker.compose.volume"] == "" {
					    continue
					}
  			  stack.Volumes = append(
  			      stack.Volumes,
  			      &proto.Volume{
  			          Name:       volume.Name,
  			          Driver:     volume.Driver,
  			          Mountpoint: volume.Mountpoint,
  			          Labels:     volume.Labels,
  			          Scope:      volume.Scope,
  			      },
  			  )
				}
    }

		for _, network := range networkGroups[stackName] {
      	stack.Networks = append(
      	    stack.Networks,
						network,
      	)
    }


    stack.Status = calculateStackStatus(stackContainers)

    response.Stacks = append(
        response.Stacks,
        stack,
    )
  }

  return response, nil

}

func calculateStackStatus(
    containers []container.Summary,
) string {

    if len(containers) == 0 {
        return "empty"
    }

    running := 0

    for _, c := range containers {
        if c.State == "running" {
            running++
        }
    }

    switch {
    case running == len(containers):
        return "running"

    case running == 0:
        return "stopped"

    default:
        return "partial"
    }
} 
