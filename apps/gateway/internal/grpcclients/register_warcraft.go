package grpcclients

import (
	"context"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/registryclient"
	characterspb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ccChecker(cc *grpc.ClientConn) error {
	if cc == nil {
		return status.Error(codes.Unavailable, "nil grpc ClientConn (downstream not connected)")
	}
	return nil
}

func RegisterWarcraft(req registryclient.Registry) {

	req.Register("warcraft.character.WarcraftCharacterService/GetCharacters", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharactersRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ccChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharacters(ctx, req.(*characterspb.GetCharactersRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/GetCharactersById", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharacterByIdRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ccChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharactersById(ctx, req.(*characterspb.GetCharacterByIdRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/GetCharactersByCharacterId", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharacterByCharacterIdRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ccChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharactersByCharacterId(ctx, req.(*characterspb.GetCharacterByCharacterIdRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/GetCharactersByUserId", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharacterByUserIdRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ccChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharactersByUserId(ctx, req.(*characterspb.GetCharacterByUserIdRequest))
		},
	})

}
