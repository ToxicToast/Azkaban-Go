// Package character contains the gRPC-based implementation of the CharacterPort interface.
// It acts as an infrastructure adapter for communication with the character service.
package character

import (
	"context"

	"github.com/ToxicToast/Azkaban-Go/apps/gateway/internal/application/warcraft/port"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	sharedpb "github.com/ToxicToast/Azkaban-Go/proto/shared"
	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft"
)

// grpcCharacterService is a gRPC-based implementation of the CharacterPort interface.
type grpcCharacterService struct {
	client characterpb.WarcraftCharacterServiceClient
}

// GetCharacters calls the remote GetCharacters RPC with optional limit, offset and deletion flag.
func (s *grpcCharacterService) GetCharacters(ctx context.Context, limit, offset *int64, withDeleted *bool) (*characterpb.GetCharactersResponse, error) {
	req := &sharedpb.ListRequest{
		Limit:       limit,
		Offset:      offset,
		WithDeleted: withDeleted,
	}
	return s.client.GetCharacters(ctx, req)
}

func (s *grpcCharacterService) GetCharactersByID(ctx context.Context, id int64, withDeleted *bool) (*characterpb.Character, error) {
	req := &sharedpb.ByIdRequest{
		Id:          id,
		WithDeleted: withDeleted,
	}
	return s.client.GetCharactersById(ctx, req)
}

func (s *grpcCharacterService) GetCharactersByCharacterID(ctx context.Context, characterID string, withDeleted *bool) (*characterpb.Character, error) {
	req := &characterpb.GetCharacterByCharacterIdRequest{
		CharacterId: characterID,
		WithDeleted: withDeleted,
	}
	return s.client.GetCharactersByCharacterId(ctx, req)
}

func (s *grpcCharacterService) GetCharactersByUserID(ctx context.Context, userID *int64, withDeleted *bool, limit, offset *int64) (*characterpb.GetCharactersResponse, error) {
	userIDWrapped := helper.WrapInt64(userID)
	req := &sharedpb.ByUserIdRequest{
		UserId:      userIDWrapped,
		WithDeleted: withDeleted,
		Limit:       limit,
		Offset:      offset,
	}
	return s.client.GetCharactersByUserId(ctx, req)
}

func (s *grpcCharacterService) GetCharactersByGuild(ctx context.Context, guild *string, withDeleted *bool, limit, offset *int64) (*characterpb.GetCharactersResponse, error) {
	req := &characterpb.GetCharacterByGuildRequest{
		Guild:       guild,
		WithDeleted: withDeleted,
		Limit:       limit,
		Offset:      offset,
	}
	return s.client.GetCharactersByGuild(ctx, req)
}

func (s *grpcCharacterService) CreateCharacter(ctx context.Context, region, realm, name string) (*characterpb.Character, error) {
	req := &characterpb.CreateCharacterRequest{
		Region: region,
		Realm:  realm,
		Name:   name,
	}
	return s.client.CreateCharacter(ctx, req)
}

func (s *grpcCharacterService) AssignCharacter(ctx context.Context, id int64, userID *int64) (*characterpb.Character, error) {
	userIDWrapped := helper.WrapInt64(userID)
	req := &characterpb.AssignCharacterRequest{
		Id:     id,
		UserId: userIDWrapped,
	}
	return s.client.AssignCharacter(ctx, req)
}

// NewGrpcCharacterService creates a new gRPC adapter for the character service.
func NewGrpcCharacterService(client characterpb.WarcraftCharacterServiceClient) port.CharacterPort {
	return &grpcCharacterService{
		client: client,
	}
}
