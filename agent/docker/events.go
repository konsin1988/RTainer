package docker

import (
    "context"

    "github.com/docker/docker/api/types/events"
    "github.com/docker/docker/api/types/filters"
)

type EventsRequest struct {
    Types   []string
    Actions []string
}

type Event struct {
    Time       int64
    Type       string
    Action     string 
    ID         string
    Attributes map[string]string
}


func (c *Client) Events(
    ctx context.Context,
    req EventsRequest,
) (<-chan Event, <-chan error) {

    f := filters.NewArgs()

    for _, t := range req.Types {
        f.Add("type", t)
    }

    for _, a := range req.Actions {
        f.Add("event", a)
    }

    dockerEvents, dockerErrors := c.cli.Events(
        ctx,
        events.ListOptions{
            Filters: f,
        },
    )

    out := make(chan Event)

    go func() {
        defer close(out)

        for e := range dockerEvents {
            out <- Event{
                Time:       e.Time,
                Type:       string(e.Type),
                Action:     string(e.Action),
                ID:         e.Actor.ID,
                Attributes: e.Actor.Attributes,
            }
        }
    }()

    return out, dockerErrors
}
