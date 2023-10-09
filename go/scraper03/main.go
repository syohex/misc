package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"
)

var start int
var count int
var released = true

func init() {
	var inclusive bool
	flag.BoolVar(&inclusive, "i", false, "start from zero")
	flag.IntVar(&count, "c", 10, "print list count")
	flag.Parse()

	if inclusive {
		start = 0
		count = count - 1
	} else {
		start = 1
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
}

var listTemplate = `|[[{{.ID}}>{{.URL}}]]|[[{{.SmallImage}}>{{.LargeImage}}]]|{{.Title}}|{{.Performer}}|{{.Date}}||`
var separator = `|~NO|PHOTO|TITLE|ACTRESS|RELEASE|NOTE|`

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
	return fmt.Sprintf("%s%03d", strings.ToLower(id), number)
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
		d.Title = e.Text
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

	var performers []string
	c.OnHTML("td span#performer a", func(e *colly.HTMLElement) {
		performers = append(performers, strings.TrimSpace(e.Text))
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
		fmt.Printf("failed to get base data %s", baseData.URL)
		return 1
	}

	if baseData.LargeImage == "" {
		fmt.Printf("failed to parse base data %s", baseData.URL)
		return 1
	}

	var b bytes.Buffer
	for i := start; i <= count; i++ {
		num := baseNumber + i

		d := &Data{
			ID:  toProductID(productID, num),
			URL: generateURL(baseURL, baseID, productID, num),
		}

		if released {
			if err := d.scrape(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return 1
			}
		}

		if d.LargeImage == "" || d.SmallImage == "" {
			d.LargeImage = generateImageURL(baseData.LargeImage, baseID, productID, num)
			d.SmallImage = generateImageURL(baseData.SmallImage, baseID, productID, num)
			d.Performer = "[[ ]]"
			d.Date = "20--"
		}

		if err := t.Execute(&b, d); err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate from template: %v\n", err)
			return 1
		}

		b.WriteRune('\n')

		if num%10 == 0 {
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
