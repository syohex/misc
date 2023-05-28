package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gocolly/colly/v2"
)

var htmlTemplate = `
<a href="{{.URL}}" target="_blank">
<img src="{{.Image}}" alt="{{.Title}}" width=800 />
</a>

<p>
</p>
`

type Data struct {
	Title string
	URL   string
	Image string
}

func copyToClipboard(text string) error {
	cmd := exec.Command("xsel", "--input", "--clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func dmm(url string, t *template.Template) (*Data, error) {
	c := colly.NewCollector()

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "age_check_done",
		Value:  "1",
		Path:   "/",
		Domain: ".dmm.co.jp",
	})

	d := &Data{
		URL: url,
	}

	if err := c.SetCookies("https://www.dmm.co.jp", cookies); err != nil {
		return nil, err
	}

	c.OnHTML("a[target]", func(e *colly.HTMLElement) {
		target := e.Attr("target")
		if d.Image == "" && target == "_package" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1[id]", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		if id == "title" {
			d.Title = e.Text
		}
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return d, nil
}

func sokmil(url string, t *template.Template) (*Data, error) {
	c := colly.NewCollector()

	d := &Data{
		URL: url,
	}

	c.OnHTML("a.sokmil-lightbox-jacket", func(e *colly.HTMLElement) {
		if d.Image == "" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1.page-title", func(e *colly.HTMLElement) {
		d.Title = e.Text
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return d, nil
}

func knights(url string, t *template.Template) (*Data, error) {
	c := colly.NewCollector()

	d := &Data{
		URL: url,
	}

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

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return d, nil
}

func mgs(url string, t *template.Template) (*Data, error) {
	c := colly.NewCollector()

	d := &Data{
		URL: url,
	}

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:   "adc",
		Value:  "1",
		Path:   "/",
		Domain: "mgstage.com",
	})

	if err := c.SetCookies("https://mgstage.com", cookies); err != nil {
		return nil, err
	}

	c.OnHTML("a.link_magnify", func(e *colly.HTMLElement) {
		if d.Image == "" {
			d.Image = e.Attr("href")
		}
	})
	c.OnHTML("h1.tag", func(e *colly.HTMLElement) {
		d.Title = strings.TrimSpace(e.Text)
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return d, nil
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

	url := os.Args[1]
	var d *Data
	if strings.Contains(url, "dmm.co.jp") {
		d, err = dmm(url, t)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(url, "sokmil.com") {
		d, err = sokmil(url, t)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(url, "knights-visual.com") {
		d, err = knights(url, t)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(url, "mgstage.com") {
		d, err = mgs(url, t)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("unsuppoerted url: %s", url)
	}

	if err := t.Execute(os.Stdout, &d); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := copyToClipboard(d.Title); err != nil {
		fmt.Printf("failed to copy title to clipboard %v\n", err)
	}

	return 0
}
