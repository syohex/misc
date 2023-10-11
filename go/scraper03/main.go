package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"
)

var ids []int
var zeros int
var released = true

func init() {
	var start int
	var end int

	flag.IntVar(&start, "s", -1, "start number")
	flag.IntVar(&end, "e", -1, "end number")
	flag.IntVar(&zeros, "z", 3, "end number")
	flag.Parse()

	if start != -1 && end == -1 {
		ids = append(ids, start)
	} else if start != -1 && end != -1 {
		for i := start; i <= end; i++ {
			ids = append(ids, i)
		}
	}
}

type Data struct {
	ID         string
	Title      string
	URL        string
	Date       string
	SmallImage string
	LargeImage string
	Performer  string
	Note       string
}

var listTemplate = `|[[{{.ID}}>{{.URL}}]]|[[{{.SmallImage}}>{{.LargeImage}}]]|{{.Title}}|{{.Performer}}|{{.Date}}|{{.Note}}|`
var separator = `|~NO|PHOTO|TITLE|ACTRESS|RELEASE|NOTE|`
var performerRegex = regexp.MustCompile(`^([^(（ ]+)`)

func formatPerformers(ps []string) string {
	if len(ps) == 0 {
		return "[[ ]]"
	}

	var ss []string
	for _, p := range ps {
		ss = append(ss, fmt.Sprintf("[[%s]]", p))
	}

	return strings.Join(ss, "／")
}

func toProductID(id string, number int) string {
	return fmt.Sprintf("%s-%03d", strings.ToUpper(id), number)
}

func idInURL(id string, number int) string {
	return fmt.Sprintf("%s%0*d", strings.ToLower(id), zeros, number)
}

func (d *Data) scrape() error {
	if strings.Contains(d.URL, "dmm.co.jp") {
		return d.dmm()
	}

	return fmt.Errorf("unsupported URL: %s", d.URL)
}

func generateURL(baseURL string, baseID string, id string, num int) string {
	return strings.ReplaceAll(baseURL, baseID, idInURL(id, num))
}

func generateImageURL(baseURL string, baseID string, id string, num int) string {
	return strings.ReplaceAll(baseURL, baseID, idInURL(id, num))
}

func stripPerformer(s string) string {
	s = strings.TrimSpace(s)
	m := performerRegex.FindStringSubmatch(s)
	if len(m) == 0 {
		return s
	}

	return m[1]
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

	c.OnHTML("h1#title", func(e *colly.HTMLElement) {
		d.Title = strings.TrimSpace(e.Text)
	})

	state := ""
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		if d.Date != "" {
			return
		}

		if d.Date == "" && (strings.HasPrefix(state, "発売日") || strings.HasPrefix(state, "商品発売日")) {
			d.Date = strings.ReplaceAll(strings.TrimSpace(e.Text), "/", "-")
			return
		}
		state = e.Text
	})

	c.OnHTML("tr td a", func(e *colly.HTMLElement) {
		if d.Note != "" {
			return
		}

		link := e.Attr("href")
		if strings.Contains(link, "id=6608") || strings.Contains(e.Text, "総集編") {
			d.Note = "総集編作品"
		}
	})

	var performers []string
	c.OnHTML("td span#performer a", func(e *colly.HTMLElement) {
		performers = append(performers, stripPerformer(e.Text))
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

	if err := c.Visit(d.URL); err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			released = false
			return nil
		}

		return err
	}

	d.Performer = formatPerformers(performers)
	return nil
}

func _main() int {
	args := flag.Args()
	if len(args) < 3 {
		fmt.Printf("Usage: %s product_id number url\n", os.Args[0])
		return 1
	}

	productID := args[0]
	numberStr := args[1]
	baseURL := args[2]

	for _, s := range args[3:] {
		id, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid argument %s\n", s)
			os.Exit(1)
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		fmt.Fprintf(os.Stderr, "ID is not specified\n")
		return 1
	}

	baseNumber, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Printf("invalid product number %s: %v", numberStr, err)
		return 1
	}

	t, err := template.New("test").Parse(listTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse template: %v\n", err)
		return 1
	}

	baseID := idInURL(productID, baseNumber)
	baseData := &Data{
		ID:  toProductID(productID, baseNumber),
		URL: generateURL(baseURL, baseID, productID, baseNumber),
	}

	if err := baseData.scrape(); err != nil {
		fmt.Printf("failed to get base data %s, %v\n", baseData.URL, err)
		return 1
	}

	if baseData.LargeImage == "" {
		fmt.Printf("failed to parse base data %s\n", baseData.URL)
		return 1
	}

	var b bytes.Buffer
	for _, id := range ids {
		d := &Data{
			ID:  toProductID(productID, id),
			URL: generateURL(baseURL, baseID, productID, id),
		}

		if released {
			if err := d.scrape(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return 1
			}
		}

		if d.LargeImage == "" || d.SmallImage == "" {
			d.LargeImage = generateImageURL(baseData.LargeImage, baseID, productID, id)
			d.SmallImage = generateImageURL(baseData.SmallImage, baseID, productID, id)
			d.Performer = "[[ ]]"
			d.Date = "20--"
		}

		if err := t.Execute(&b, d); err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate from template: %v\n", err)
			return 1
		}

		b.WriteRune('\n')

		if id%10 == 0 {
			b.WriteString(separator)
			b.WriteRune('\n')
		}
	}

	output := b.String()
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
