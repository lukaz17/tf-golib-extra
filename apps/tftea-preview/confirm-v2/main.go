package main

import (
	"fmt"
	"os"

	tftea "github.com/tforce-io/tf-golib-extra/tftea-v2"
)

func main() {
	yes, err := tftea.NewConfirm().
		WithLabel("Do you like Go?").
		WithValue(true).
		Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("Thanks for confirm: %v", yes)
}
