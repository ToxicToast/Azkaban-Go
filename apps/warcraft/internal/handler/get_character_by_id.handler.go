package handler

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetCharacterByIdHandler(ctx context.Context, in *characterpb.GetCharacterByIdRequest) (*characterpb.Character, error) {
	characters := BuildCharacters()
	for _, character := range characters {
		if character.Id == in.GetId() {
			return character, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "character not found")
}
