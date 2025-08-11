//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func pusherBlog() error {
	fmt.Println("🔍 Dockerize Blog Service...")
	cmd := exec.Command("docker", "push", "toxictoast/blog-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pusherCronjob() error {
	fmt.Println("🔍 Dockerize Cronjob Service...")
	cmd := exec.Command("docker", "push", "toxictoast/cronjob-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pusherFoodfolio() error {
	fmt.Println("🔍 Dockerize Foodfolio Service...")
	cmd := exec.Command("docker", "push", "toxictoast/foodfolio-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pusherGateway() error {
	fmt.Println("🔍 Dockerize Gateway Service...")
	cmd := exec.Command("docker", "push", "toxictoast/gateway-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pusherTwitch() error {
	fmt.Println("🔍 Dockerize Twitch Service...")
	cmd := exec.Command("docker", "push", "toxictoast/twitch-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pusherWarcraft() error {
	fmt.Println("🔍 Dockerize Warcraft Service...")
	cmd := exec.Command("docker", "push", "toxictoast/warcraft-go:dev")
	fmt.Printf("→ %s\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Push() error {
	fmt.Println("🔍 Dockerize Services...")
	var errPushBlog = pusherBlog()
	if errPushBlog != nil {
		return fmt.Errorf("failed to dockerize blog: %w", errPushBlog)
	}

	var errPushCronjob = pusherCronjob()
	if errPushCronjob != nil {
		return fmt.Errorf("failed to dockerize cronjob: %w", errPushCronjob)
	}

	var errPushFoodfolio = pusherFoodfolio()
	if errPushFoodfolio != nil {
		return fmt.Errorf("failed to dockerize foodfolio: %w", errPushFoodfolio)
	}

	var errPushGateway = pusherGateway()
	if errPushGateway != nil {
		return fmt.Errorf("failed to dockerize gateway: %w", errPushGateway)
	}

	var errPushTwitch = pusherTwitch()
	if errPushTwitch != nil {
		return fmt.Errorf("failed to dockerize twitch: %w", errPushTwitch)
	}

	var errPushWarcraft = pusherWarcraft()
	if errPushWarcraft != nil {
		return fmt.Errorf("failed to dockerize warcraft: %w", errPushWarcraft)
	}

	return nil
}
