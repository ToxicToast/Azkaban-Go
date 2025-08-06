// Package port defines application-level interfaces (ports) for the character domain.
// These ports are implemented by infrastructure adapters like gRPC or REST clients.
package port

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/libs/shared/proto/warcraft"
)

// CharacterPort defines the application interface for accessing and modifying characters.
// It is implemented by infrastructure adapters like gRPC or REST.
type CharacterPort interface {
	GetCharacters(ctx context.Context, limit, offset *int64, withDeleted *bool) (*characterpb.GetCharactersResponse, error)
	GetCharactersByID(ctx context.Context, id int64, withDeleted *bool) (*characterpb.Character, error)
	GetCharactersByCharacterID(ctx context.Context, characterID string, withDeleted *bool) (*characterpb.Character, error)
	GetCharactersByUserID(ctx context.Context, userID *int64, withDeleted *bool, limit, offset *int64) (*characterpb.GetCharactersResponse, error)
	GetCharactersByGuild(ctx context.Context, guild *string, withDeleted *bool, limit, offset *int64) (*characterpb.GetCharactersResponse, error)
	CreateCharacter(ctx context.Context, region, realm, name string) (*characterpb.Character, error)
	AssignCharacter(ctx context.Context, id int64, userID *int64) (*characterpb.Character, error)
}
