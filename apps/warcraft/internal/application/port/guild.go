package port

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/libs/shared/proto/warcraft"
)

// GuildPort defines the application interface for accessing and modifying guilds.
// It is implemented by infrastructure adapters like gRPC or REST.
type GuildPort interface {
	GetGuilds(ctx context.Context, limit, offset *int64, withDeleted *bool) (*characterpb.GetGuildsResponse, error)
}