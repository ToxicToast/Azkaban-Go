package grpc

import (
	"context"
	"log"
	"net"

	"github.com/ToxicToast/Azkaban-Go/apps/warcraft/internal/handler"
	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type svc struct {
	characterpb.UnimplementedWarcraftCharacterServiceServer
}

type Client struct {
	server   *grpc.Server
	health   *health.Server
	listener net.Listener
}

func (s *svc) GetCharacters(ctx context.Context, in *characterpb.GetCharactersRequest) (*characterpb.GetCharactersResponse, error) {
	return handler.GetCharactersHandler(ctx, in)
}

func (s *svc) GetCharactersById(ctx context.Context, in *characterpb.GetCharacterByIdRequest) (*characterpb.Character, error) {
	return handler.GetCharacterByIdHandler(ctx, in)
}

func (s *svc) GetCharactersByCharacterId(ctx context.Context, in *characterpb.GetCharacterByCharacterIdRequest) (*characterpb.Character, error) {
	return handler.GetCharacterByCharacterIdHandler(ctx, in)
}

func (s *svc) CreateCharacter(ctx context.Context, in *characterpb.CreateCharacterRequest) (*characterpb.Character, error) {
	return handler.CreateCharacterHandler(ctx, in)
}

func (s *svc) AssignCharacter(ctx context.Context, in *characterpb.AssignCharacterRequest) (*characterpb.Character, error) {
	return nil, nil
}

func NewGRPCHandler(port string) *Client {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	server := grpc.NewServer()
	hs := health.NewServer()

	return &Client{
		listener: lis,
		server:   server,
		health:   hs,
	}
}

func (h *Client) Register() {
	characterpb.RegisterWarcraftCharacterServiceServer(h.server, &svc{})
	healthpb.RegisterHealthServer(h.server, h.health)
	h.health.SetServingStatus("warcraft", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(h.server)
}

func (h *Client) Serve() {
	if err := h.server.Serve(h.listener); err != nil {
		panic("failed to serve: " + err.Error())
	}
	log.Println("Warcraft Service listening on :8082")
}
