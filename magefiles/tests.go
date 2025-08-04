//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func Test() error {
	fmt.Println("✅ Running all tests...")

	cmd := exec.Command("go", "test", "-v", "./_test")

	fmt.Printf("→ %s\n", cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
