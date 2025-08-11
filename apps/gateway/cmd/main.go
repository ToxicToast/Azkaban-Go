package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ToxicToast/Azkaban-Go/apps/gateway/internal/grpcclients"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/configclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/grpcclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/registryclient"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/routerclient"
)

var root = "."

var (
	cfg        *configclient.Config
	router     *routerclient.Client
	registry   registryclient.Registry
	grpcClient *grpcclient.Client
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
	router.BuildHealthRoute([]string{"warcraft"}, cfg.Routes, registry)
}

func init() {
	loadConfigs()
	loadRegistry()
	loadGrpcClients()
	loadRouter()
}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting %s (v%s) in %s mode...\n", cfg.App.Name, cfg.App.Version, cfg.App.Env)
	fmt.Printf("On Port %v...\n", cfg.Server.Http.Port)

	go router.RunServer()
	<-stop
	fmt.Println("Shutting down gateway service...")
}
