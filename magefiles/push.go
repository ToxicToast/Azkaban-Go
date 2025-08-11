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
	cmd := exec.Command("docker", "push", "-t", "toxictoast/blog-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerCronjob() error {
	fmt.Println("🔍 Dockerize Cronjob Service...")
	cmd := exec.Command("docker", "push", "-t", "toxictoast/cronjob-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerFoodfolio() error {
	fmt.Println("🔍 Dockerize Foodfolio Service...")
	cmd := exec.Command("docker", "push", "-t", "toxictoast/foodfolio-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerGateway() error {
	fmt.Println("🔍 Dockerize Gateway Service...")
	cmd := exec.Command("docker", "push", "-t", "toxictoast/gateway-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerTwitch() error {
	fmt.Println("🔍 Dockerize Twitch Service...")
	cmd := exec.Command("docker", "push", "-t", "toxictoast/twitch-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func dockerWarcraft() error {
	fmt.Println("🔍 Dockerize Warcraft Service...")
	cmd := exec.Command("docker", "push", "-t", "toxictoast/warcraft-go:dev")
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
