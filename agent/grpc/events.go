package grpcserver

import (
    pb "konsin1988/agent/proto"
)

func (s *Server) Events(
    req *pb.EventsRequest,
    stream pb.AgentService_EventsServer,
) error {
    eventsCh, errCh := s.containerSvc.Events(
        stream.Context(),
        req,
    )

    for {
        select {

        case e, ok := <-eventsCh:
            if !ok {
                return nil
            }

            if err := stream.Send(&pb.EventMessage{
                Time:       e.Time,
                Type:       e.Type,
                Action:     e.Action,
                Id:         e.ID,
                Attributes: e.Attributes,
            }); err != nil {
                return err
            }

        case err := <-errCh:
            if err != nil {
                return err
            }

        case <-stream.Context().Done():
            return nil
        }
    }
}
