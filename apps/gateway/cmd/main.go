package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ToxicToast/Azkaban-Go/apps/gateway/internal/grpcclients"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/configclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/grpcclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/healthmon"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/registryclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/routerclient"
)

var root = "."

var (
	bgCtx      context.Context
	cfg        *configclient.Config
	router     *routerclient.Client
	registry   registryclient.Registry
	grpcClient *grpcclient.Client
	monitor    *healthmon.Monitor
)

func loadConfigs() {
	configclient.LoadEnvFile(root)
	envName := configclient.GetenvDefault("APP_ENV", "dev")
	configPaths := []string{
		filepath.Join("configs", envName, "app.yaml"),
		filepath.Join("configs", envName, "auth.yaml"),
		filepath.Join("configs", envName, "cache.yaml"),
		filepath.Join("configs", envName, "downstreams.yaml"),
		filepath.Join("configs", envName, "events.yaml"),
		filepath.Join("configs", envName, "features.yaml"),
		filepath.Join("configs", envName, "health.yaml"),
		filepath.Join("configs", envName, "observability.yaml"),
		filepath.Join("configs", envName, "server.yaml"),
		filepath.Join("configs", envName, "routes.yaml"),
	}
	cfgClient := configclient.NewClient(root, configPaths)
	cfg = cfgClient.Load()
}

func loadRegistry() {
	req := registryclient.NewRegistry()
	grpcclients.RegisterWarcraft(req)
	registry = req
}

func loadGrpcClients() {
	grpcClient = grpcclient.NewClient(cfg.Downstreams.Services)
}

func loadRouter() {
	router = routerclient.NewClient(cfg.App.Env, ":"+cfg.Server.Http.Port, grpcClient)
	router.BuildRoutes(cfg.Routes, registry)
	router.BuildHealthRoute(monitor, cfg.Health.Services, cfg.Routes, registry, cfg.Health.LivenessPath, cfg.Health.ReadinessPath)
}

func loadMonitor() {

	monitor = healthmon.NewMonitor(
		cfg.Health.Services,
		func(c context.Context, svc string) error {
			return grpcClient.Ping(c, svc)
		},
		cfg.Health.CheckIntervals.GrpcTargets,
		cfg.Health.CheckIntervals.Redis,
		cfg.Health.CheckIntervals.Kafka,
	)
	monitor.Start(bgCtx)
}

func init() {
	bgCtx = context.Background()
	loadConfigs()
	loadRegistry()
	loadGrpcClients()
	loadMonitor()
	loadRouter()
}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go router.RunServer()
	<-stop
	fmt.Println("Shutting down gateway service...")
}
