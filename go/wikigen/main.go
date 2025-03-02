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
	"github.com/gocolly/colly/v2"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
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

type Data struct {
	ID         string
	Title      string
	Date       string
	SmallImage string
	LargeImage string
	Config     Config
}

func (d *Data) scrape(productURL string) error {
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

func (d *Data) dmm(productURL string) error {
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

func (d *Data) sokmil(productURL string) error {
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

func validateArgument(url1, url2 string) bool {
	return strings.Contains(url1, "www.sokmil.com") && strings.Contains(url2, "www.dmm.co.jp")
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

func (d *Data) String(url1 string, url2 string, config *Config) (string, error) {
	sokmilAff, err := sokmilAffiliateURL(url1, config)
	if err != nil {
		return "", err
	}
	dmmAff, err := dmmAffiliateURL(url2, config)
	if err != nil {
		return "", err
	}

	output := ""
	// id part
	output += "|"
	output += fmt.Sprintf("[[%s>%s]]", d.ID, dmmAff)

	// image part
	output += "|"
	output += fmt.Sprintf("center:[[&ref(%s,180)>%s]]", d.SmallImage, d.LargeImage)
	output += "~~"
	output += fmt.Sprintf("[[ソクミル>%s]] [[FANZA>%s]]", sokmilAff, dmmAff)

	// title part
	output += "|"
	output += d.Title

	// performer part
	output += "|"

	// date part
	output += "|"
	output += d.Date

	// end part
	output += "||"

	return output, nil
}

func _main() int {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s ID url1 url2\n", os.Args[0])
		return 1
	}

	id := strings.ToUpper(os.Args[1])
	url1 := os.Args[2]
	url2 := os.Args[3]

	if !validateArgument(url1, url2) {
		fmt.Fprintf(os.Stderr, "url1 should be sokmil, url2 should be fanza(url1=%s, url2=%s)\n", url1, url2)
		return 1
	}

	config, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config file: %v", err)
		return 1
	}

	d := &Data{
		ID: id,
	}
	if err := d.scrape(url1); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	output, err := d.String(url1, url2, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert data to string: %v\n", err)
		return 1
	}

	fmt.Println(output)

	if err := clipboard.WriteAll(output); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy text into clipboard: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(_main())
}
