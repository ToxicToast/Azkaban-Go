package grpcclients

import (
	"context"
	"fmt"

	"github.com/ToxicToast/Azkaban-Go/libs/shared/registryclient"
	characterspb "github.com/ToxicToast/Azkaban-Go/proto/warcraft/character"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func RegisterWarcraft(req registryclient.Registry) {

	req.Register("warcraft.character.WarcraftCharacterService/GetCharacters", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharactersRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ConnectionChecker(cc)
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
			err := ConnectionChecker(cc)
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
			err := ConnectionChecker(cc)
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
			err := ConnectionChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharactersByUserId(ctx, req.(*characterspb.GetCharacterByUserIdRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/GetCharactersByGuild", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.GetCharacterByGuildRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ConnectionChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.GetCharactersByGuild(ctx, req.(*characterspb.GetCharacterByGuildRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/CreateCharacter", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.CreateCharacterRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ConnectionChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			fmt.Printf("Creating character with request: %+v\n", req)
			return c.CreateCharacter(ctx, req.(*characterspb.CreateCharacterRequest))
		},
	})

	req.Register("warcraft.character.WarcraftCharacterService/AssignCharacter", registryclient.Entry{
		NewReq: func() proto.Message { return &characterspb.AssignCharacterRequest{} },
		Invoke: func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error) {
			err := ConnectionChecker(cc)
			if err != nil {
				return nil, err
			}
			c := characterspb.NewWarcraftCharacterServiceClient(cc)
			return c.AssignCharacter(ctx, req.(*characterspb.AssignCharacterRequest))
		},
	})

}
