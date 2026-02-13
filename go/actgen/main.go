package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	Name         string              `yaml:"name"`
	Image        string              `yaml:"image"`
	Aliases      map[string]string   `yaml:"aliases"`
	SNS          map[string]string   `yaml:"sns"`
	Fanza        string              `yaml:"fanza"`
	Sokmil       string              `yaml:"sokmil"`
	RelatedPages map[string][]string `yaml:"related_pages"`
}

func sortMap(m map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	values := make([]string, 0, len(m))
	for _, k := range keys {
		values = append(values, m[k])
	}

	return keys, values
}

func (a *Actress) Render(conf *Config) (string, error) {
	var err error
	hasFanzaUrl := a.Fanza != ""
	hasSokmilUrl := a.Sokmil != ""

	a.Fanza, err = dmmAffiliateURL(a.Fanza, conf)
	if err != nil {
		return "", err
	}

	a.Sokmil, err = sokmilAffiliateURL(a.Sokmil, conf)
	if err != nil {
		return "", err
	}

	for k, v := range a.Aliases {
		if v == "" {
			a.Aliases[k] = ""
		} else {
			affiliateURL, err := dmmAffiliateURL(v, conf)
			if err != nil {
				return "", err
			}
			a.Aliases[k] = affiliateURL
		}
	}

	var sb strings.Builder
	if a.Image != "" {
		sb.WriteString(fmt.Sprintf("&ref(%s)\n", a.Image))
		sb.WriteString("\n")
	}

	if len(a.Aliases) > 0 {
		sb.WriteString("** 別名\n")
		for name, url := range a.Aliases {
			if url == "" {
				sb.WriteString(fmt.Sprintf("- %s\n", name))
			} else {
				sb.WriteString(fmt.Sprintf("- [[%s>%s]]\n", name, url))
			}
		}
		sb.WriteString("\n")
	}

	if hasFanzaUrl || hasSokmilUrl {
		sb.WriteString("** 作品リンク\n")
		if hasFanzaUrl {
			sb.WriteString(fmt.Sprintf("- [[FANZA>%s]]\n", a.Fanza))
		}
		if hasSokmilUrl {
			sb.WriteString(fmt.Sprintf("- [[ソクミル>%s]]\n", a.Sokmil))
		}
	}
	sb.WriteString("\n")

	if len(a.SNS) > 0 {
		sb.WriteString("** SNS\n")

		names, urls := sortMap(a.SNS)
		for i := range names {
			if urls[i] != "" {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", names[i], urls[i]))
			}
		}
		sb.WriteString("\n")
	}

	if len(a.RelatedPages) > 0 {
		var makers []string
		for k := range a.RelatedPages {
			makers = append(makers, k)
		}

		slices.Sort(makers)

		sb.WriteString("** 関連ページ\n")
		for _, maker := range makers {
			sb.WriteString(fmt.Sprintf("- %s\n", maker))

			products := a.RelatedPages[maker]
			slices.Sort(products)

			for _, product := range products {
				sb.WriteString(fmt.Sprintf("-- [[%s]]\n", product))
			}
		}
	}

	return sb.String(), nil
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
