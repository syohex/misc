package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/syohex/clipboard"
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
	Image        string            `yaml:"image"`
	Aliases      map[string]string `yaml:"aliases"`
	SNS          map[string]string `yaml:"sns"`
	Instagram    string            `yaml:"instagram"`
	Tiktok       string            `yaml:"tictok"`
	Fanza        string            `yaml:"fanza"`
	Sokmil       string            `yaml:"sokmil"`
	RelatedPages []string          `yaml:"related_pages"`
}

var pageTemplate = `&ref({{.Image}})

** 別名
{{range $name, $url := .Aliases}}- [[{{$name}}>{{$url}}]]}
{{end}}
** 作品リンク
- [[FANZA>{{.Fanza}}]]
- [[Sokmil>{{.Sokmil}}]]

** SNS
{{range $name, $val := .SNS }}- {{ $name }}: {{ $val }}
{{end}}
** 関連ページ
{{range .RelatedPages}}- [[{{.}}]]
{{end}}`

func (a *Actress) Render(conf *Config) (string, error) {
	var err error

	a.Fanza, err = dmmAffiliateURL(a.Fanza, conf)
	if err != nil {
		return "", err
	}

	a.Sokmil, err = sokmilAffiliateURL(a.Sokmil, conf)
	if err != nil {
		return "", err
	}

	for k, v := range a.Aliases {
		affiliateURL, err := dmmAffiliateURL(v, conf)
		if err != nil {
			return "", err
		}
		a.Aliases[k] = affiliateURL
	}

	t, err := template.New("page").Parse(pageTemplate)
	if err != nil {
		return "", err
	}

	b := bytes.NewBufferString("")
	err = t.Execute(b, a)
	if err != nil {
		return "", err
	}

	return b.String(), nil
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

	output, err := actress.Render(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to render: %v\n", err)
		return 1
	}

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
