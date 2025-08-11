package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		lis, err := net.Listen("tcp", ":8083")
		if err != nil {
			log.Fatal(err)
		}

		server := grpc.NewServer()

		hs := health.NewServer()
		healthpb.RegisterHealthServer(server, hs)
		hs.SetServingStatus("twitch", healthpb.HealthCheckResponse_SERVING)

		reflection.Register(server)

		log.Println("Twitch stub listening on :8083")
		log.Fatal(server.Serve(lis))
	}()

	<-stop
	fmt.Println("Shutting down twitch service...")
}
