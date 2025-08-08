// Package ports defines application-level interfaces (ports) for a specific domain.
// These ports are implemented by http adapters like REST clients.
package ports

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft"
)

// GuildPort defines the interface for accessing and modifying guilds.
// It is implemented by http adapters like REST.
type GuildPort interface {
	GetGuilds(ctx context.Context, limit, offset *int64, withDeleted *bool) (*characterpb.GetGuildsResponse, error)
}
