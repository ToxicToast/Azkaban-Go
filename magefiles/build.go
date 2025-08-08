//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildBlog() error {
	fmt.Println("🔍 Build Blog Service...")
	cmd := exec.Command("go", "build", "-o", "bin/blog", "apps/blog/cmd/main.go")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildCronjob() error {
	fmt.Println("🔍 Build Cronjob Service...")
	cmd := exec.Command("go", "build", "-o", "bin/cronjob", "apps/cronjob/cmd/main.go")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildGateway() error {
	fmt.Println("🔍 Build Gateway Service...")
	cmd := exec.Command("go", "build", "-o", "bin/gateway", "apps/gateway/cmd/main.go")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildTwitch() error {
	fmt.Println("🔍 Build Twitch Service...")
	cmd := exec.Command("go", "build", "-o", "bin/twitch", "apps/twitch/cmd/main.go")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildWarcraft() error {
	fmt.Println("🔍 Build Warcraft Service...")
	cmd := exec.Command("go", "build", "-o", "bin/warcraft", "apps/warcraft/cmd/main.go")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Build() error {
	fmt.Println("🔍 Build Services...")

	var errBlog = buildBlog()
	if errBlog != nil {
		return fmt.Errorf("failed to build blog: %w", errBlog)
	}

	var errCronjob = buildCronjob()
	if errCronjob != nil {
		return fmt.Errorf("failed to build cronjob: %w", errCronjob)
	}

	var errGateway = buildGateway()
	if errGateway != nil {
		return fmt.Errorf("failed to build gateway: %w", errGateway)
	}

	var errTwitch = buildTwitch()
	if errTwitch != nil {
		return fmt.Errorf("failed to build twitch: %w", errTwitch)
	}

	var errWarcraft = buildWarcraft()
	if errWarcraft != nil {
		return fmt.Errorf("failed to build warcraft: %w", errTwitch)
	}

	return nil
}
