//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func dockerBlog() error {
	fmt.Println("🔍 Dockerize Blog Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/blog-go:dev", "--build-arg", "BUILD_PATH=./apps/blog/cmd/main.go", "--build-arg", "BINARY_NAME=blog", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerCronjob() error {
	fmt.Println("🔍 Dockerize Cronjob Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/cronjob-go:dev", "--build-arg", "BUILD_PATH=./apps/cronjob/cmd/main.go", "--build-arg", "BINARY_NAME=cronjob", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerFoodfolio() error {
	fmt.Println("🔍 Dockerize Foodfolio Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/foodfolio-go:dev", "--build-arg", "BUILD_PATH=./apps/foodfolio/cmd/main.go", "--build-arg", "BINARY_NAME=foodfolio", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerGateway() error {
	fmt.Println("🔍 Dockerize Gateway Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/gateway-go:dev", "--build-arg", "BUILD_PATH=./apps/gateway/cmd/main.go", "--build-arg", "BINARY_NAME=gateway", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerTwitch() error {
	fmt.Println("🔍 Dockerize Twitch Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/twitch-go:dev", "--build-arg", "BUILD_PATH=./apps/twitch/cmd/main.go", "--build-arg", "BINARY_NAME=twitch", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerWarcraft() error {
	fmt.Println("🔍 Dockerize Warcraft Service...")
	cmd := exec.Command("docker", "build", "-t", "toxictoast/warcraft-go:dev", "--build-arg", "BUILD_PATH=./apps/warcraft/cmd/main.go", "--build-arg", "BINARY_NAME=warcraft", "--build-arg", "VERSION=0.0.1", ".")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Docker() error {
	fmt.Println("🔍 Dockerize Services...")
	var errBlog = dockerBlog()
	if errBlog != nil {
		return fmt.Errorf("failed to dockerize blog: %w", errBlog)
	}

	var errCronjob = dockerCronjob()
	if errCronjob != nil {
		return fmt.Errorf("failed to dockerize cronjob: %w", errCronjob)
	}

	var errFoodfolio = dockerFoodfolio()
	if errFoodfolio != nil {
		return fmt.Errorf("failed to dockerize foodfolio: %w", errFoodfolio)
	}

	var errGateway = dockerGateway()
	if errGateway != nil {
		return fmt.Errorf("failed to dockerize gateway: %w", errGateway)
	}

	var errTwitch = dockerTwitch()
	if errTwitch != nil {
		return fmt.Errorf("failed to dockerize twitch: %w", errTwitch)
	}

	var errWarcraft = dockerWarcraft()
	if errWarcraft != nil {
		return fmt.Errorf("failed to dockerize warcraft: %w", errWarcraft)
	}

	return nil
}
