package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/syohex/clipboard"
	"github.com/goccy/go-yaml"
)

type Config struct {
	Dmm    AffiliateInfo `yaml:"dmm"`
	Sokmil AffiliateInfo `yaml:"sokmil"`
}

type AffiliateInfo struct {
	Id string `yaml:"id"`
}

func readConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(dir, ".config", "blog", "config.yaml")
	bs, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(bs, conf)
	return conf, err
}

type Actress struct {
	Name         string            `yaml:"name"`
	Aliases      []string          `yaml:"aliases"`
	SNS          map[string]string `yaml:"sns"`
	Instagram    string            `yaml:"instagram"`
	Tiktok       string            `yaml:"tictok"`
	Products     map[string]string `yaml:"products"`
	RelatedPages []string          `yaml:"related_pages"`
}

func (a *Actress) Render(conf *Config) {

}

func sokmilAffiliateURL(productURL string, config *Config) (string, error) {
	u, err := url.Parse(productURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("affi", config.Sokmil.Id)
	q.Set("utm_source", "sokmil_ad")
	q.Set("utm_medium", "affiliate")
	q.Set("utm_campaign", config.Sokmil.Id)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func dmmAffiliateURL(productURL string, config *Config) (string, error) {
	u, err := url.Parse("https://al.dmm.co.jp/")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("lurl", productURL)
	q.Set("af_id", config.Dmm.Id)
	q.Set("ch", "link_tool")
	q.Set("ch_id", "link")

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func _main() int {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s data.yaml\n", os.Args[0])
		return 1
	}

	inputFile := os.Args[1]
	c, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input yaml file: %v\n", err)
		return 1
	}

	var actress Actress
	if err := yaml.Unmarshal(c, &actress); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input yaml file(%s): %v\n", inputFile, err)
		return 1
	}

	config, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config file: %v\n", err)
		return 1
	}
	fmt.Println(config)

	var sb strings.Builder
	output := sb.String()
	fmt.Print(output)

	if err := clipboard.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy text into clipboard: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(_main())
}
