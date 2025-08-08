//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runTests() error {
	fmt.Println("✅ Running all tests...")

	cmd := exec.Command("go", "test", "-v", "-race", "./_test")

	fmt.Printf("→ %s\n", cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runCoverage() error {
	fmt.Println("✅ Running test coverage...")

	cmd := exec.Command("go", "tool", "cover", "-html=coverage.out")

	fmt.Printf("→ %s\n", cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Test() error {
	fmt.Println("✅ Running all test steps...")

	var errTest = runTests()
	if errTest != nil {
		return fmt.Errorf("failed to run tests: %w", errTest)
	}

	var errCov = runCoverage()
	if errCov != nil {
		return fmt.Errorf("failed to run tests: %w", errCov)
	}

	return nil
}
