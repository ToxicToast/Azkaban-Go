package guild

import (
	"context"

	"github.com/ToxicToast/Azkaban-Go/apps/gateway/internal/application/warcraft/port"
	sharedpb "github.com/ToxicToast/Azkaban-Go/proto/shared"
	guildpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft"
)

// grpcGuildService is a gRPC-based implementation of the GuildPort interface.
type grpcGuildService struct {
	client guildpb.WarcraftGuildServiceClient
}

// GetGuilds calls the remote GetGuilds RPC with optional limit, offset and deletion flag.
func (s *grpcGuildService) GetGuilds(ctx context.Context, limit, offset *int64, withDeleted *bool) (*guildpb.GetGuildsResponse, error) {
	req := &sharedpb.ListRequest{
		Limit:       limit,
		Offset:      offset,
		WithDeleted: withDeleted,
	}
	return s.client.GetGuilds(ctx, req)
}

func NewGrpcGuildService(client guildpb.WarcraftGuildServiceClient) port.GuildPort {
	return &grpcGuildService{
		client: client,
	}
}
