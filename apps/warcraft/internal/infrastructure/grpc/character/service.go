// Package character contains the gRPC-based implementation of the CharacterPort interface.
// It acts as an infrastructure adapter for communication with the character service.
package character

import (
	"context"

	"github.com/ToxicToast/Azkaban-Go/apps/warcraft/internal/application/port"
	"github.com/ToxicToast/Azkaban-Go/libs/shared/helper"
	characterpb "github.com/ToxicToast/Azkaban-Go/libs/shared/proto/warcraft"
)

// grpcCharacterService is a gRPC-based implementation of the CharacterPort interface.
type grpcCharacterService struct {
	client characterpb.WarcraftCharacterServiceClient
}

// GetCharacters calls the remote GetCharacters RPC with optional limit, offset and deletion flag.
func (s *grpcCharacterService) GetCharacters(ctx context.Context, limit, offset *int64, withDeleted *bool) (*characterpb.GetCharactersResponse, error) {
	req := &characterpb.GetCharactersRequest{
		Limit:       limit,
		Offset:      offset,
		WithDeleted: withDeleted,
	}
	return s.client.GetCharacters(ctx, req)
}

func (s *grpcCharacterService) GetCharactersByID(ctx context.Context, id int64, withDeleted *bool) (*characterpb.Character, error) {
	req := &characterpb.GetCharacterByIdRequest{
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
	req := &characterpb.GetCharacterByUserIdRequest{
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

func (s *grpcCharacterService) UpdateCharacter(ctx context.Context, id int64, displayRealm, displayName, gender, faction, race, class, spec *string, level, itemLevel *int64, loggedinAt *string) (*characterpb.Character, error) {
	levelWrapped := helper.WrapInt64(level)
	itemLevelWrapped := helper.WrapInt64(itemLevel)
	req := &characterpb.UpdateCharacterRequest{
		Id:           id,
		DisplayRealm: displayRealm,
		DisplayName:  displayName,
		Gender:       gender,
		Faction:      faction,
		Race:         race,
		Class:        class,
		Spec:         spec,
		Level:        levelWrapped,
		ItemLevel:    itemLevelWrapped,
		LoggedinAt:   loggedinAt,
	}
	return s.client.UpdateCharacter(ctx, req)
}

func (s *grpcCharacterService) UpdateCharacterGuild(ctx context.Context, id int64, guild *string, rank *int64, oldGuild *string) (*characterpb.Character, error) {
	rankWrapped := helper.WrapInt64(rank)
	req := &characterpb.UpdateCharacterGuildRequest{
		Id:       id,
		Guild:    guild,
		Rank:     rankWrapped,
		OldGuild: oldGuild,
	}
	return s.client.UpdateCharacterGuild(ctx, req)
}

func (s *grpcCharacterService) UpdateCharacterMythic(ctx context.Context, id int64, mythic *int64) (*characterpb.Character, error) {
	mythicWrapped := helper.WrapInt64(mythic)
	req := &characterpb.UpdateCharacterMythicRequest{
		Id:     id,
		Mythic: mythicWrapped,
	}
	return s.client.UpdateCharacterMythic(ctx, req)
}

func (s *grpcCharacterService) UpdateCharacterRaid(ctx context.Context, id int64, raid *string) (*characterpb.Character, error) {
	req := &characterpb.UpdateCharacterRaidRequest{
		Id:   id,
		Raid: raid,
	}
	return s.client.UpdateCharacterRaid(ctx, req)
}

func (s *grpcCharacterService) UpdateCharacterMedia(ctx context.Context, id int64, inset, avatar *string) (*characterpb.Character, error) {
	req := &characterpb.UpdateCharacterMediaRequest{
		Id:     id,
		Inset:  inset,
		Avatar: avatar,
	}
	return s.client.UpdateCharacterMedia(ctx, req)
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
