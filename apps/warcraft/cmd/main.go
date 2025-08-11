package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	grpc2 "github.com/ToxicToast/Azkaban-Go/apps/warcraft/internal/grpc"
)

var grpcHandler *grpc2.Client

func init() {
	grpcHandler = grpc2.NewGRPCHandler("8082")
	grpcHandler.Register()
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go grpcHandler.Serve()

	<-stop
	fmt.Println("Shutting down warcraft service...")

}
