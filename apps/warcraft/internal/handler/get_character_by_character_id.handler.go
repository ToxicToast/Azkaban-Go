package handler

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetCharacterByCharacterIdHandler(ctx context.Context, in *characterpb.GetCharacterByCharacterIdRequest) (*characterpb.Character, error) {
	characters := BuildCharacters()
	for _, character := range characters {
		if character.CharacterId == in.GetCharacterId() {
			return character, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "character not found")
}
