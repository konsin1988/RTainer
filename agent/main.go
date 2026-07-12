package main

import (
    "log"
    "net"

    "konsin1988/agent/config"
    "konsin1988/agent/docker"
		"konsin1988/agent/service"
    grpcserver "konsin1988/agent/grpc"

    "google.golang.org/grpc"
		"google.golang.org/grpc/health"
		healthpb "google.golang.org/grpc/health/grpc_health_v1"
		"google.golang.org/grpc/reflection"
)

func main() {

    cfg := config.Load()

    dockerClient, err := docker.New()
    if err != nil {
        log.Fatal(err)
    }

		containerSvc := service.NewContainerService(dockerClient);
		imageSvc := service.NewImageService(dockerClient);
		networkSvc := service.NewNetworkService(dockerClient);
		handler := grpcserver.New(containerSvc, imageSvc, networkSvc);

    listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
    if err != nil {
        log.Fatal(err)
    }

		grpcSrv := grpc.NewServer();
		handler.Register(grpcSrv);

		// health 
		healthServer := health.NewServer();
		healthServer.SetServingStatus(
			"",
			healthpb.HealthCheckResponse_SERVING,
		);
		healthpb.RegisterHealthServer(
			grpcSrv,
			healthServer,
		)

		reflection.Register(grpcSrv);

    log.Printf(
        "agent listening on :%s",
        cfg.GRPCPort,
    )

    if err :=
        grpcSrv.Serve(listener); err != nil {

        log.Fatal(err)
    }
}
