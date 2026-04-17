package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tforce-io/tf-golib-extra/tftea-v2"
)

func main() {
	countries := []*tftea.SelectOption{
		{Key: "au", Label: "Australia"},
		{Key: "br", Label: "Brazil"},
		{Key: "ca", Label: "Canada"},
		{Key: "cn", Label: "China"},
		{Key: "de", Label: "Germany"},
		{Key: "eg", Label: "Egypt"},
		{Key: "fr", Label: "France"},
		{Key: "gb", Label: "United Kingdom"},
		{Key: "in", Label: "India"},
		{Key: "id", Label: "Indonesia"},
		{Key: "it", Label: "Italy"},
		{Key: "jp", Label: "Japan"},
		{Key: "mx", Label: "Mexico"},
		{Key: "ng", Label: "Nigeria"},
		{Key: "nz", Label: "New Zealand"},
		{Key: "ru", Label: "Russia"},
		{Key: "sa", Label: "Saudi Arabia"},
		{Key: "za", Label: "South Africa"},
		{Key: "kr", Label: "South Korea"},
		{Key: "es", Label: "Spain"},
		{Key: "se", Label: "Sweden"},
		{Key: "tr", Label: "Turkey"},
		{Key: "us", Label: "United States"},
		{Key: "vn", Label: "Vietnam"},
	}
	chosen, err := tftea.NewSelect().
		WithLabel("Select a country:").
		WithOptions(countries).
		Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("Selected country: %v\n\n", chosen[0])

	languages := []*tftea.SelectOption{
		{Key: "go", Label: "Go"},
		{Key: "ts", Label: "TypeScript"},
		{Key: "py", Label: "Python"},
		{Key: "rs", Label: "Rust"},
		{Key: "rb", Label: "Ruby"},
		{Key: "java", Label: "Java"},
		{Key: "cs", Label: "C#"},
		{Key: "cpp", Label: "C++"},
		{Key: "kt", Label: "Kotlin"},
		{Key: "sw", Label: "Swift"},
		{Key: "php", Label: "PHP"},
		{Key: "ex", Label: "Elixir"},
		{Key: "hs", Label: "Haskell"},
		{Key: "scala", Label: "Scala"},
		{Key: "dart", Label: "Dart"},
	}
	prefixFilter := func(filter string, all []*tftea.SelectOption) []*tftea.SelectOption {
		if filter == "" {
			return all
		}
		lower := strings.ToLower(filter)
		out := make([]*tftea.SelectOption, 0)
		for _, o := range all {
			if strings.HasPrefix(strings.ToLower(o.Label), lower) {
				out = append(out, o)
			}
		}
		return out
	}
	picks, err := tftea.NewSelect().
		WithLabel("Select your favourite languages:").
		WithMultiSelect(true).
		WithOptions(languages).
		WithFilter(prefixFilter).
		Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("Selected languages: %v\n\n", picks)
}
