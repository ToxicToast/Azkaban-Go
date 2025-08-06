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
	folders = []string{
		// "apps/blog",
		// "apps/cronjob",
		// "apps/foodfolio",
		// "apps/gateway",
		// "apps/twitch",
		"apps/warcraft",
		// "libs/auth",
		// "libs/events",
		"libs/shared",
	}
)

func LintFolders() error {
	fmt.Println("🔍 Linting selected folders...")
	config := ".golangci.yml"

	for _, dir := range folders {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			fmt.Printf("→ Skipping %s (no go.mod)\n", dir)
			continue
		}

		hasGoFiles := false
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".go" {
				hasGoFiles = true
				return filepath.SkipDir // abort early
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to scan %s: %w", dir, err)
		}

		if !hasGoFiles {
			fmt.Printf("→ Skipping %s (no Go files)\n", dir)
			continue
		}

		fmt.Printf("→ Linting %s\n", dir)
		cmd := exec.Command("golangci-lint", "run", "--config", config, "--timeout=3m", dir+"/...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("lint failed in %s: %w", dir, err)
		}
	}

	fmt.Println("✅ Linting finished.")
	return nil
}

func Lint() error {
	return LintFolders()
}
