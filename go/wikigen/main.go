package main

import (
	"bytes"
	"embed"
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
var makerLabelTable map[string]string
var withoutLabelLink bool
var isSeries = false

func init() {
	flag.StringVar(&overrideLabel, "l", "", "label override with given value")
	flag.BoolVar(&isSeries, "s", false, "use series instead of label")
	flag.BoolVar(&withoutLabelLink, "n", false, "without label link")
	flag.Parse()

	initMakerTable()
}

//go:embed maker.list
var listFile embed.FS

func initMakerTable() {
	makerLabelTable = make(map[string]string)

	data, err := listFile.ReadFile("maker.list")
	if err != nil {
		panic("cannot read maker.list")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		columns := strings.SplitN(line, "\t", 2)
		makerLabelTable[columns[0]] = columns[1]
	}
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
	Size       string
	ListName   string
}

var wikiTemplate = `//{{.Date}} {{.ID}}
[[{{.Title}}（{{.MakerLabel}}）>{{.URL}}]]　[[({{.ListName}}一覧)>{{.Label}}]]
[[{{.SmallImage}}>{{.LargeImage}}]]`

var wikiTemplateWithoutLabel = `//{{.Date}} {{.ID}}
[[{{.Title}}（{{.MakerLabel}}）>{{.URL}}]]
[[{{.SmallImage}}>{{.LargeImage}}]]`

var videoCTemplate = `//{{.Date}} {{.ID}}
[[{{.Title}}{{.Size}}>{{.URL}}]]　[[({{.ListName}}一覧)>{{.Label}}]]
[[&ref({{.SmallImage}},147)>{{.SmallImage}}]]`

var videoCTemplateWithoutLabel = `//{{.Date}} {{.ID}}
[[{{.Title}}{{.Size}}>{{.URL}}]]
[[&ref({{.SmallImage}},147)>{{.SmallImage}}]]`

var idRegex = regexp.MustCompile(`([a-zA-Z]+)(\d+)(?:so|z)?$`)
var labelRegex = regexp.MustCompile(`^([^(（]+)`)
var performerRegex = regexp.MustCompile(`^([^(（]+)`)

func isSameMeaning(maker string, label string) bool {
	for k, v := range makerLabelTable {
		if k == maker && v == label {
			return true
		}
	}
	return false
}

func makerLabel(maker string, label string) string {
	if maker == label || strings.Contains(maker, label) {
		return maker
	}
	if isSameMeaning(maker, label) {
		return maker
	}

	return fmt.Sprintf("%s／%s", maker, label)
}

func convertID(id string) string {
	m := idRegex.FindStringSubmatch(id)
	if m == nil {
		return id
	}

	charPart := m[1]
	numPart := m[2]
	if len(numPart) == 5 { // e.g. 00123
		numPart = strings.TrimPrefix(numPart, "00")
		if len(numPart) == 5 { // e.g 01234
			numPart = strings.TrimPrefix(numPart, "0")
		}
	} else if len(numPart) == 4 { // e.g. 0001
		numPart = strings.TrimPrefix(numPart, "0")
	}

	if charPart == "dsvr" && !strings.HasPrefix(id, "dsvr") {
		charPart = "3dsvr"
	}

	return fmt.Sprintf("%s-%s", strings.ToUpper(charPart), numPart)
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

	return strings.TrimSpace(m[1])
}

func stripPerformer(s string) string {
	s = strings.TrimSpace(s)
	m := performerRegex.FindStringSubmatch(s)
	if len(m) == 0 {
		return s
	}

	return strings.TrimSpace(m[1])
}

func packageImage(img string) string {
	if !strings.HasSuffix(img, "pl.jpg") {
		return img
	}

	return strings.ReplaceAll(img, "pl.jpg", "ps.jpg")
}

var dandyRegex = regexp.MustCompile(`1dandy(\d+)`)

func extractDandyID(url string) string {
	m := dandyRegex.FindStringSubmatch(url)
	if m == nil {
		return ""
	}

	return strings.TrimSpace(m[1])
}

var titleReplacer = strings.NewReplacer("@", "＠")

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
		d.Title = titleReplacer.Replace(strings.TrimSpace(e.Text))
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

	c.OnHTML("td > a", func(e *colly.HTMLElement) {
		if d.Maker != "" && d.Label != "" {
			return
		}

		link := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		if strings.Contains(link, "=maker") || strings.Contains(link, "maker=") {
			d.Maker = normalizeLabel(text)
		} else if strings.Contains(link, "=label") || strings.Contains(link, "label=") {
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

		d.SmallImage = packageImage(e.Attr("content"))
	})

	if err := c.Visit(url); err != nil {
		return err
	}

	if d.SmallImage != "" && d.LargeImage == "" {
		d.LargeImage = strings.Replace(d.SmallImage, "ps.jpg", "pl.jpg", 1)
	}
	if strings.Contains(d.URL, "1dandy") {
		if id := extractDandyID(d.URL); id != "" {
			tmpl := "https://pics.dmm.co.jp/digital/video/1dandy###/1dandy###jp-1.jpg"
			d.LargeImage = strings.ReplaceAll(tmpl, "###", id)
		}
	}

	return nil
}

var agePartRe = regexp.MustCompile(`\([0-9]+\)`)

func combineTitleAndName(name string, title string) string {
	m := agePartRe.FindStringSubmatch(name)
	if len(m) == 0 {
		return title
	}

	return title + m[0]
}

func (d *Data) dmmTypeC(url string) error {
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

	title := ""
	c.OnHTML("h1#title", func(e *colly.HTMLElement) {
		title = strings.TrimSpace(e.Text)
	})

	state := ""
	name := ""
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		if d.Date != "" && d.ID != "" && d.Size != "" && name != "" {
			return
		}

		text := strings.TrimSpace(e.Text)
		if d.Date == "" && strings.HasPrefix(state, "配信開始日") {
			d.Date = strings.ReplaceAll(strings.TrimSpace(text), "/", ".")
			return
		} else if d.ID == "" && strings.HasPrefix(state, "品番") {
			d.ID = convertID(text)
			return
		} else if d.Size == "" && strings.HasPrefix(state, "サイズ") {
			d.Size = strings.TrimSpace(e.Text)
			return
		} else if name == "" && strings.HasPrefix(state, "名前") {
			name = strings.TrimSpace(e.Text)
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
		if strings.Contains(link, "=maker") || strings.Contains(link, "maker=") {
			d.Maker = text
		} else if strings.Contains(link, "label=") || strings.Contains(link, "maker=") {
			d.Label = normalizeLabel(text)
		}
	})

	c.OnHTML("td span#performer a", func(e *colly.HTMLElement) {
		d.Performers = append(d.Performers, e.Text)
	})

	c.OnHTML("meta[property=og\\:image]", func(e *colly.HTMLElement) {
		if d.SmallImage != "" {
			return
		}

		d.SmallImage = packageImage(e.Attr("content"))
	})

	c.OnHTML("a[name=package-image]", func(e *colly.HTMLElement) {
		if d.LargeImage != "" {
			return
		}

		d.LargeImage = e.Attr("href")
	})

	if err := c.Visit(url); err != nil {
		return err
	}

	if name != "" && title != "" {
		d.Title = combineTitleAndName(name, title)
	} else if name != "" {
		d.Title = name
	} else {
		d.Title = title
	}
	d.Title = titleReplacer.Replace(d.Title)

	if strings.Contains(d.Size, "---") {
		d.Size = ""
	} else {
		d.Size = " " + d.Size
	}

	if d.Label != "" {
		d.Title = fmt.Sprintf("%s %s", d.Label, d.Title)
	}

	return nil
}

func _main() int {
	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("Usage: %s URL\n", os.Args[0])
		return 1
	}

	url := args[0]

	var templateStr string
	if strings.Contains(url, "dmm.co.jp") && strings.Contains(url, "videoc") {
		if withoutLabelLink {
			templateStr = videoCTemplateWithoutLabel
		} else {
			templateStr = videoCTemplate
		}
	} else {
		if withoutLabelLink {
			templateStr = wikiTemplateWithoutLabel
		} else {
			templateStr = wikiTemplate
		}
	}

	t, err := template.New("test").Parse(templateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse template: %v\n", err)
		return 1
	}

	d := &Data{URL: url}

	if strings.Contains(url, "dmm.co.jp") {
		if strings.Contains(url, "videoc") {
			err = d.dmmTypeC(url)
		} else {
			err = d.dmm(url)
		}
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

	if isSeries {
		d.ListName = "シリーズ"
	} else {
		d.ListName = "レーベル"
	}

	var b bytes.Buffer
	if err := t.Execute(&b, d); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate template: %v\n", err)
		return 1
	}

	output := b.String()
	if len(d.Performers) > 1 {
		output += fmt.Sprintf("\n出演者：%s", formatPerformers(d.Performers))
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
