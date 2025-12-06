package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/syohex/clipboard"
)

type Actress struct {
	Name         string              `yaml:"name"`
	Image        string              `yaml:"image"`
	Aliases      map[string]string   `yaml:"aliases"`
	SNS          map[string]string   `yaml:"sns"`
	Fanza        string              `yaml:"fanza"`
	Sokmil       string              `yaml:"sokmil"`
	RelatedPages map[string][]string `yaml:"related_pages"`
}

func _main() int {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s data.yaml\n", os.Args[0])
		return 1
	}

	inputDir := os.Args[1]
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read directory %s: %v\n", inputDir, err)
		return 1
	}

	var names []string
	for _, entry := range entries {
		yamlPath := filepath.Join(inputDir, entry.Name())

		f, err := os.ReadFile(yamlPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file %s: %v\n", yamlPath, err)
			return 1
		}

		var actress Actress
		if err := yaml.Unmarshal(f, &actress); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal yaml file %s: %v\n", yamlPath, err)
			return 1
		}

		names = append(names, actress.Name)
	}

	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	var sb strings.Builder
	for _, name := range names {
		sb.WriteString(fmt.Sprintf("- [[%s]]\n", name))
	}

	output := sb.String()
	if err := clipboard.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy output into clipboard: %v\n", err)
		return 1
	}

	fmt.Print(output)
	return 0
}

func main() {
	os.Exit(_main())
}
