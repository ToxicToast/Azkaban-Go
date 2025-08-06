//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	protoDir   = "libs/proto"
	outputDir  = "libs/shared/proto"
	protoFiles = []string{
		"shared/common.proto",
		//
		"foodfolio/category.proto",
		"foodfolio/company.proto",
		"foodfolio/item.proto",
		"foodfolio/location.proto",
		"foodfolio/size.proto",
		"foodfolio/type.proto",
		"foodfolio/warehouse.proto",
		//
		"warcraft/common.proto",
		//
		"warcraft/types/character-type.proto",
		"warcraft/types/character-requests.proto",
		"warcraft/types/character-responses.proto",
		"warcraft/types/guild-type.proto",
		"warcraft/types/guild-responses.proto",
		//
		"warcraft/character.proto",
		"warcraft/guild.proto",
		//
		"warcraft/cronjobs/cronjob-search.proto",
		"warcraft/cronjobs/cronjob-guild.proto",
		"warcraft/cronjobs/cronjob-guild-roster.proto",
		"warcraft/cronjobs/cronjob-character-profile.proto",
		"warcraft/cronjobs/cronjob-media.proto",
		"warcraft/cronjobs/cronjob-mythic.proto",
		//
		"warcraft/cronjob.proto",
		//
		"twitch/viewer.proto",
	}
)

func Protobuf() error {
	fmt.Println("🔧 Generating Go code from .proto files...")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	for _, file := range protoFiles {
		fullPath := filepath.ToSlash(filepath.Join(protoDir, file))
		fmt.Printf("→ %s\n", fullPath)
		cmd := exec.Command("protoc",
			"--go_out="+outputDir,
			"--go_opt=paths=source_relative",
			"--go-grpc_out="+outputDir,
			"--go-grpc_opt=paths=source_relative",
			"--proto_path="+protoDir,
			fullPath,
		)
		fmt.Printf("→ %s\n", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("protoc failed for %s: %w", file, err)
		}
	}

	fmt.Println("✅ Protobuf generation complete.")
	return nil
}
