package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/goccy/go-yaml"
	"github.com/gocolly/colly/v2"
)

const (
	userAgent           = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
	tableHeader         = "|~ID|Image|タイトル|出演者(出演順)|発売日|Note|"
	actressSeparator    = "／"
	wikiNewLine         = "~~"
	wikiColumnSeparator = '|'
)

type Config struct {
	Dmm    AffiliateInfo `json:"dmm"`
	Sokmil AffiliateInfo `json:"sokmil"`
}

type AffiliateInfo struct {
	Id string `json:"id"`
}

func readConfig() (*Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(dir, ".config", "blog", "config.json")
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = json.NewDecoder(file).Decode(conf)
	return conf, err
}

type ProductData struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	SokmilURL string   `yaml:"sokmil"`
	FanzaURL  string   `yaml:"fanza"`
	Actresses []string `yaml:"actresses"`
	Note      string   `yaml:"note"`
}

type Product struct {
	ID         string
	Title      string
	Date       string
	SmallImage string
	LargeImage string
	Config     Config
	Actresses  []string
	Note       string
}

func (d *Product) scrape(productURL string) error {
	if strings.Contains(productURL, "www.sokmil.com") {
		return d.sokmil(productURL)
	} else if strings.Contains(productURL, "dmm.co.jp") {
		return d.dmm(productURL)
	}

	return fmt.Errorf("unsupported URL: %s", productURL)
}

func pageImage(s string) string {
	if !strings.HasSuffix(s, "pl.jpg") {
		return s
	}

	return strings.Replace(s, "pl.jpg", "ps.jpg", 1)
}

func isDateState(s string) bool {
	states := []string{"発売日", "商品発売日"}
	for _, state := range states {
		if strings.HasPrefix(s, state) {
			return true
		}
	}

	return false
}

var titleReplacer = strings.NewReplacer("@", "＠")

func formatDate(dateStr string) string {
	return strings.ReplaceAll(strings.TrimSpace(dateStr), "/", "-")
}

func (d *Product) dmm(productURL string) error {
	c := colly.NewCollector()
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "age_check_done",
		Value:  "1",
		Path:   "/",
		Domain: ".dmm.co.jp",
	})

	if err := c.SetCookies("https://www.dmm.co.jp", cookies); err != nil {
		return err
	}

	c.OnHTML("h1#title", func(e *colly.HTMLElement) {
		d.Title = titleReplacer.Replace(strings.TrimSpace(e.Text))
	})

	state := ""
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		if d.Date == "" && isDateState(state) {
			d.Date = formatDate(e.Text)
			return
		}

		state = e.Text
	})

	c.OnHTML("meta[property=og\\:image]", func(e *colly.HTMLElement) {
		if d.SmallImage != "" {
			return
		}

		d.SmallImage = pageImage(e.Attr("content"))
	})

	if err := c.Visit(productURL); err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil
		}

		return err
	}

	if d.SmallImage != "" && d.LargeImage == "" {
		d.LargeImage = strings.Replace(d.SmallImage, "ps.jpg", "pl.jpg", 1)
	}

	return nil
}

func (d *Product) sokmil(productURL string) error {
	c := colly.NewCollector(
		colly.UserAgent(userAgent),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	})

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "AGEAUTH",
		Value:  "ok",
		Path:   "/",
		Domain: ".sokmil.com",
	})

	if err := c.SetCookies("https://www.sokmil.com", cookies); err != nil {
		return err
	}

	c.OnHTML("a.sokmil-lightbox-jacket", func(e *colly.HTMLElement) {
		if d.LargeImage == "" {
			d.LargeImage = e.Attr("href")
		}
	})
	c.OnHTML("img.jacket-img", func(e *colly.HTMLElement) {
		if d.SmallImage == "" {
			d.SmallImage = e.Attr("content")
		}
	})
	c.OnHTML("span[itemprop=releaseDate]", func(e *colly.HTMLElement) {
		if d.Date == "" {
			d.Date = formatDate(e.Attr("content"))
		}
	})
	c.OnHTML("h1.page-title", func(e *colly.HTMLElement) {
		d.Title = e.Text
	})
	return c.Visit(productURL)
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

func (d *Product) String(url1 string, url2 string, config *Config) (string, error) {
	sokmilAff, err := sokmilAffiliateURL(url1, config)
	if err != nil {
		return "", err
	}
	dmmAff, err := dmmAffiliateURL(url2, config)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// id part
	sb.WriteRune(wikiColumnSeparator)
	sb.WriteString(fmt.Sprintf("[[%s>%s]]", d.ID, dmmAff))

	// image part
	sb.WriteRune(wikiColumnSeparator)
	sb.WriteString(fmt.Sprintf("center:[[&ref(%s,180)>%s]]", d.SmallImage, d.LargeImage))
	sb.WriteString(wikiNewLine)
	sb.WriteString(fmt.Sprintf("[[ソクミル>%s]] [[FANZA>%s]]", sokmilAff, dmmAff))

	// title part
	sb.WriteRune(wikiColumnSeparator)
	sb.WriteString(d.Title)

	// performer part
	sb.WriteRune(wikiColumnSeparator)

	var actStrs []string
	for _, actress := range d.Actresses {
		if strings.HasSuffix(actress, "?") {
			actStrs = append(actStrs, strings.TrimRight(actress, "?"))
		} else {
			actStrs = append(actStrs, fmt.Sprintf("[[%s]]", actress))
		}
	}
	sb.WriteString(strings.Join(actStrs, actressSeparator))

	// release date part
	sb.WriteRune(wikiColumnSeparator)
	sb.WriteString(d.Date)

	// note part
	sb.WriteRune(wikiColumnSeparator)
	sb.WriteString(d.Note)
	sb.WriteRune(wikiColumnSeparator)

	return sb.String(), nil
}

func matchID(productID string, ids []string) bool {
	if len(ids) == 0 {
		return true
	}

	for _, id := range ids {
		if productID == id {
			return true
		}
	}

	return false
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

	var filterIDs []string
	if len(os.Args) > 2 {
		for _, id := range os.Args[2:] {
			filterIDs = append(filterIDs, strings.ToUpper(id))
		}
	}

	var productData []ProductData
	if err := yaml.Unmarshal(c, &productData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input yaml file(%s): %v\n", inputFile, err)
		return 1
	}

	config, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config file: %v\n", err)
		return 1
	}

	var sb strings.Builder
	if len(filterIDs) == 0 {
		sb.WriteString(tableHeader)
		sb.WriteRune('\n')
	}

	for _, product := range productData {
		pd := &Product{
			ID:        product.ID,
			Actresses: product.Actresses,
			Note:      product.Note,
		}

		if !matchID(pd.ID, filterIDs) {
			continue
		}

		if err := pd.scrape(product.SokmilURL); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		output, err := pd.String(product.SokmilURL, product.FanzaURL, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert data to string: %v\n", err)
			return 1
		}

		sb.WriteString(output)
		sb.WriteRune('\n')
	}

	output := sb.String()
	fmt.Print(output)

	if err := clipboard.WriteAll(output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy text into clipboard: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(_main())
}
