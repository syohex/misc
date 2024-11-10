package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

var htmlTemplate = `
<a href="{{.URL}}" target="_blank">
<img src="{{.Image}}" alt="{{.Title}}" {{.Width}} />
</a>

<p>
</p>
`

type Config struct {
	Dmm    AffiliateInfo `json:"dmm"`
	Sokmil AffiliateInfo `json:"sokmil"`
}

type AffiliateInfo struct {
	Id string `json:"id"`
}

type Data struct {
	Title  string
	URL    string
	Image  string
	Width  template.HTMLAttr
	config Config
}

func (d *Data) readConfig() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(dir, ".config", "blog", "config.json")
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(&d.config)
}

func (d *Data) dmm() error {
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

	// video a
	c.OnHTML("a[target]", func(e *colly.HTMLElement) {
		target := e.Attr("target")
		if d.Image == "" && target == "_package" {
			d.Image = e.Attr("href")
		}
	})
	// mono
	c.OnHTML("a[name]", func(e *colly.HTMLElement) {
		name := e.Attr("name")
		if d.Image == "" && name == "package-image" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1[id]", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		if id == "title" {
			d.Title = e.Text
		}
	})

	u, err := url.Parse("https://al.dmm.co.jp/")
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("lurl", d.URL)
	q.Set("af_id", d.config.Dmm.Id)
	q.Set("ch", "link_tool")
	q.Set("ch_id", "link")

	u.RawQuery = q.Encode()
	d.URL = u.String()

	return c.Visit(d.URL)
}

func (d *Data) sokmil() error {
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
		if d.Image == "" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1.page-title", func(e *colly.HTMLElement) {
		d.Title = e.Text
	})

	u, err := url.Parse(d.URL)
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("affi", d.config.Sokmil.Id)
	q.Set("utm_source", "sokmil_ad")
	q.Set("utm_medium", "affiliate")
	q.Set("utm_campaign", d.config.Sokmil.Id)

	u.RawQuery = q.Encode()
	d.URL = u.String()
	return c.Visit(d.URL)
}

func (d *Data) knights() error {
	d.Width = "width=800"

	c := colly.NewCollector()

	c.OnHTML(".entry-inner > p > a", func(e *colly.HTMLElement) {
		if d.Image == "" {
			d.Image = e.Attr("href")
			if strings.HasPrefix(d.Image, "/") {
				d.Image = "https://www.knights-visual.com" + d.Image
			}
		}
	})
	c.OnHTML("h1.entry-title", func(e *colly.HTMLElement) {
		d.Title = e.Text
	})

	return c.Visit(d.URL)
}

func (d *Data) mgs() error {
	c := colly.NewCollector()

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "adc",
		Value:  "1",
		Path:   "/",
		Domain: "mgstage.com",
	})

	if err := c.SetCookies("https://mgstage.com", cookies); err != nil {
		return err
	}

	c.OnHTML("a.link_magnify", func(e *colly.HTMLElement) {
		if d.Image == "" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1.tag", func(e *colly.HTMLElement) {
		d.Title = strings.TrimSpace(e.Text)
	})

	return c.Visit(d.URL)
}

func stripQueryString(urlStr string) (string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path), nil
}

func copyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

func main() {
	os.Exit(_main())
}

func _main() int {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s url\n", os.Args[0])
		return 1
	}

	t, err := template.New("test").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	url, err := stripQueryString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	d := &Data{
		URL: url,
	}

	if err := d.readConfig(); err != nil {
		log.Fatal(err)
	}

	if strings.Contains(url, "dmm.co.jp") {
		err = d.dmm()
	} else if strings.Contains(url, "sokmil.com") {
		err = d.sokmil()
	} else if strings.Contains(url, "knights-visual.com") {
		err = d.knights()
	} else if strings.Contains(url, "mgstage.com") {
		err = d.mgs()
	} else {
		log.Fatalf("unsuppoerted url: %s", url)
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(os.Stdout, d); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := copyToClipboard(d.Title); err != nil {
		fmt.Printf("failed to copy title to clipboard %v\n", err)
	}

	return 0
}
