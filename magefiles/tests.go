//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

var testFolder = "_test"

func testForFiles() bool {
	fmt.Println("🔍 Checking for test files...")

	cmd := exec.Command("find", "./"+testFolder, "-type", "f", "-name", "*_test.go")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("No test files found.")
		return false
	}

	fmt.Println("✅ Running all tests...")
	return true
}

func runTests() error {
	if testForFiles() {
		cmd := exec.Command("go", "test", "-v", "-race", "./"+testFolder)
		fmt.Printf("→ %s\n", cmd)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
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
