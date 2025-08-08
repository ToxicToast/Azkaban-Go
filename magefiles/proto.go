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
	outputDir  = "proto"
	protoFiles = []string{
		// "foodfolio/category.proto",
		// "foodfolio/company.proto",
		// "foodfolio/item.proto",
		// "foodfolio/location.proto",
		// "foodfolio/size.proto",
		// "foodfolio/type.proto",
		// "foodfolio/warehouse.proto",
		//
		"warcraft/common.proto",
		"warcraft/character/character-type.proto",
		"warcraft/character/character-requests.proto",
		"warcraft/character/character-responses.proto",
		// "warcraft/guild/guild-type.proto",
		// "warcraft/guild/guild-responses.proto",
		"warcraft/character/character.proto",
		// "warcraft/guild/guild.proto",
		// "warcraft/cronjob/cronjob-search.proto",
		// "warcraft/cronjob/cronjob-guild.proto",
		// "warcraft/cronjob/cronjob-guild-roster.proto",
		// "warcraft/cronjob/cronjob-character-profile.proto",
		// "warcraft/cronjob/cronjob-media.proto",
		// "warcraft/cronjob/cronjob-mythic.proto",
		// "warcraft/cronjob/cronjob.proto",
		//
		// "twitch/viewer.proto",
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
