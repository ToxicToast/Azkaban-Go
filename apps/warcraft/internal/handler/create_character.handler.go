package handler

import (
	"context"

	characterpb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
)

func CreateCharacterHandler(ctx context.Context, in *characterpb.CreateCharacterRequest) (*characterpb.Character, error) {
	character := &characterpb.Character{
		Id:     0,
		Region: in.GetRegion(),
		Realm:  in.GetRealm(),
		Name:   in.GetName(),
	}
	return character, nil
}
