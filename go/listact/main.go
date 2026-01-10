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

var dirMap map[string][]string
var kanaMap map[string][]string

func init() {
	dirMap = map[string][]string{
		"a":  {"a", "i", "u", "e", "o"},
		"ka": {"ka", "ki", "ku", "ke", "ko"},
		"sa": {"sa", "shi", "su", "se", "so"},
		"ta": {"ta", "chi", "tsu", "te", "to"},
		"na": {"na", "ni", "nu", "ne", "no"},
		"ha": {"ha", "hi", "fu", "he", "ho"},
		"ma": {"ma", "mi", "mu", "me", "mo"},
		"ya": {"ya", "yu", "yo"},
		"ra": {"ra", "ri", "ru", "re", "ro"},
		"wa": {"wa"},
	}

	kanaMap = map[string][]string{
		"a":  {"あ", "い", "う", "え", "お"},
		"ka": {"か", "き", "く", "け", "こ"},
		"sa": {"さ", "し", "す", "せ", "そ"},
		"ta": {"た", "ち", "つ", "て", "と"},
		"na": {"な", "に", "ぬ", "ね", "の"},
		"ha": {"は", "ひ", "ふ", "へ", "ほ"},
		"ma": {"ま", "み", "む", "め", "も"},
		"ya": {"や", "ゆ", "よ"},
		"ra": {"ら", "り", "る", "れ", "ろ"},
		"wa": {"わ"},
	}
}

func _main() int {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s data.yaml\n", os.Args[0])
		return 1
	}

	input := os.Args[1]
	if _, ok := dirMap[input]; !ok {
		fmt.Fprintf(os.Stderr, "invalid input: %s. Input must be a, ka, sa, ta, na, ha, ma, ya, ra, wa\n", input)
		return 1
	}

	var sb strings.Builder
	kanas := kanaMap[input]
	for i, dir := range dirMap[input] {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read directory %s: %v\n", dir, err)
			return 1
		}

		var actresses []struct {
			Name    string
			Aliases []string
		}
		for _, entry := range entries {
			yamlPath := filepath.Join(dir, entry.Name())

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

			aliases := []string{}
			for alias := range actress.Aliases {
				aliases = append(aliases, alias)
			}
			sort.Slice(aliases, func(i, j int) bool {
				return aliases[i] < aliases[j]
			})

			actresses = append(actresses, struct {
				Name    string
				Aliases []string
			}{
				Name:    actress.Name,
				Aliases: aliases,
			})
		}

		if len(actresses) == 0 {
			continue
		}

		sort.Slice(actresses, func(i, j int) bool {
			return actresses[i].Name < actresses[j].Name
		})

		sb.WriteString(fmt.Sprintf("**%s\n", kanas[i]))
		for _, actress := range actresses {
			if len(actress.Aliases) == 0 {
				sb.WriteString(fmt.Sprintf("- [[%s]]\n", actress.Name))
			} else {
				aliasList := strings.Join(actress.Aliases, ", ")
				sb.WriteString(fmt.Sprintf("- [[%s]] (別名義 %s)\n", actress.Name, aliasList))
			}
		}

		if i != len(dirMap[input])-1 {
			sb.WriteString("\n")
		}
	}

	output := sb.String()
	if err := clipboard.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy output into clipboard: %v\n", err)
		return 1
	}

	fmt.Println(output)
	return 0
}

func main() {
	os.Exit(_main())
}
