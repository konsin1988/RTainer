package service

import (
	"context"
	"fmt"

	"konsin1988/agent/docker"
	"konsin1988/agent/proto"
		
	"github.com/docker/docker/api/types/container"
)

type StackService struct {
  docker *docker.Client
}

func NewStackService(d *docker.Client) *StackService {
  return &StackService{docker: d}
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
    

	for _, c := range containers {
		stackName := c.Labels["com.docker.compose.project"]
		fmt.Printf("Stackname = %s", stackName)

    if stackName == "" {
        continue
    }

    groups[stackName] = append(groups[stackName], c)
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
