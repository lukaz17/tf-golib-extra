package main

import (
	"fmt"
	"os"
	"strings"

	tftea "github.com/tforce-io/tf-golib-extra/tftea-v2"
)

func main() {
	name, err := tftea.NewPrompt().
		WithLabel("Please enter your name:").
		WithPlaceholder("John Doe").
		WithValidation(ValidateRequired).
		Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("Greeting, %s!\n", name)
}

func ValidateRequired(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("input is empty")
	}
	return nil
}
