package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"
)

var overrideLabel string

func init() {
	flag.StringVar(&overrideLabel, "l", "", "label override with given value")
	flag.Parse()
}

type Data struct {
	Title      string
	URL        string
	Date       string
	Maker      string
	Label      string
	SmallImage string
	LargeImage string
	ID         string
	MakerLabel string
	Performers []string
}

var wikiTemplate = `//{{.Date}} {{.ID}}
[[{{.Title}}（{{.MakerLabel}}）>{{.URL}}]] [[（レーベル一覧）>{{.Label}}]]
[[{{.SmallImage}}>{{.LargeImage}}]]`

var idRegex = regexp.MustCompile(`([a-zA-Z]+)(\d+)$`)
var labelRegex = regexp.MustCompile(`^([^(（]+)`)
var performerRegex = regexp.MustCompile(`^([^(（]+)`)

func makerLabel(maker string, label string) string {
	if maker == label || strings.Contains(maker, label) {
		return maker
	}

	return fmt.Sprintf("%s/%s", maker, label)
}

func convertID(id string) string {
	m := idRegex.FindStringSubmatch(id)
	if m == nil {
		return id
	}

	return fmt.Sprintf("%s-%s", strings.ToUpper(m[1]), m[2])
}

func formatPerformers(ps []string) string {
	var ss []string
	for _, p := range ps {
		ss = append(ss, fmt.Sprintf("[[%s]]", p))
	}

	return strings.Join(ss, "／")
}

func normalizeLabel(label string) string {
	m := labelRegex.FindStringSubmatch(label)
	if m == nil {
		return label
	}

	return m[1]
}

func stripPerformer(s string) string {
	s = strings.TrimSpace(s)
	m := performerRegex.FindStringSubmatch(s)
	if len(m) == 0 {
		return s
	}

	return m[1]
}

func (d *Data) dmm(url string) error {
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
		d.Title = strings.TrimSpace(e.Text)
	})

	state := ""
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		if d.Date != "" && d.ID != "" {
			return
		}

		text := strings.TrimSpace(e.Text)
		if d.Date == "" && (strings.HasPrefix(state, "発売日") || strings.HasPrefix(state, "商品発売日")) {
			d.Date = strings.ReplaceAll(strings.TrimSpace(text), "/", ".")
			return
		} else if d.ID == "" && strings.HasPrefix(state, "品番") {
			d.ID = convertID(text)
			return
		}
		state = text
	})

	c.OnHTML("td a", func(e *colly.HTMLElement) {
		if d.Maker != "" && d.Label != "" {
			return
		}

		link := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		if strings.Contains(link, "maker") {
			d.Maker = text
		} else if strings.Contains(link, "label") {
			d.Label = normalizeLabel(text)
		}
	})

	c.OnHTML("td span#performer a", func(e *colly.HTMLElement) {
		d.Performers = append(d.Performers, stripPerformer(e.Text))
	})

	c.OnHTML("meta[property=og\\:image]", func(e *colly.HTMLElement) {
		if d.SmallImage != "" {
			return
		}

		d.SmallImage = e.Attr("content")
	})

	c.OnHTML("a[name=package-image]", func(e *colly.HTMLElement) {
		if d.LargeImage != "" {
			return
		}

		d.LargeImage = e.Attr("href")
	})

	return c.Visit(url)
}

func _main() int {
	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("Usage: %s URL\n", os.Args[0])
		return 1
	}

	t, err := template.New("test").Parse(wikiTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse template: %v\n", err)
		return 1
	}

	url := args[0]
	d := &Data{URL: url}

	if strings.Contains(url, "dmm.co.jp") {
		err = d.dmm(url)
	} else {
		fmt.Fprintf(os.Stderr, "unsupported URL: %s\n", url)
		return 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to extract info from %s(%v)\n", url, err)
		return 1
	}

	if d.Label == "" {
		d.Label = d.Maker
	}
	d.MakerLabel = makerLabel(d.Maker, d.Label)

	if overrideLabel != "" {
		d.Label = overrideLabel
	}

	var b bytes.Buffer
	if err := t.Execute(&b, d); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate template: %v\n", err)
		return 1
	}

	output := b.String()
	if len(d.Performers) > 1 {
		output += fmt.Sprintf("\n出演者: %s", formatPerformers(d.Performers))
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
