package handler

import (
	"context"

	"github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
)

func GetCharactersHandler(ctx context.Context, in *character.GetCharactersRequest) (*character.GetCharactersResponse, error) {
	characters := BuildCharacters()
	return &character.GetCharactersResponse{
		Data:  characters,
		Total: int64(len(characters)),
	}, nil
}