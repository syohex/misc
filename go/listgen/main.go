package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"
)

var ids []int
var zeros int
var withDirector bool
var noHeader bool
var listCount int

var listTemplate string
var baseTemplate = `|[[{{.ID}}>{{.URL}}]]|[[{{.SmallImage}}>{{.LargeImage}}]]|{{.Title}}|{{.Performer}}|{{.Date}}|{{.Note}}|`
var templateWithDirector = `|[[{{.ID}}>{{.URL}}]]|[[{{.SmallImage}}>{{.LargeImage}}]]|{{.Title}}|{{.Performer}}|{{.Director}}|{{.Date}}|{{.Note}}|`

var separator string
var baseSeparator = `|~NO|PHOTO|TITLE|ACTRESS|RELEASE|NOTE|`
var separatorWithDirector = `|~NO|PHOTO|TITLE|ACTRESS|DIRECTOR|RELEASE|NOTE|`

func init() {
	var start int
	var end int

	flag.IntVar(&start, "s", -1, "start number")
	flag.IntVar(&end, "e", -1, "end number")
	flag.IntVar(&zeros, "z", 3, "zero padding length")
	flag.BoolVar(&withDirector, "d", false, "with director column")
	flag.BoolVar(&noHeader, "n", false, "no header")
	flag.IntVar(&listCount, "c", -1, "list count")
	flag.Parse()

	if start != -1 && end == -1 {
		ids = append(ids, start)
	} else if start != -1 && end != -1 {
		for i := start; i <= end; i++ {
			ids = append(ids, i)
		}
	}

	if withDirector {
		listTemplate = templateWithDirector
		separator = separatorWithDirector
	} else {
		listTemplate = baseTemplate
		separator = baseSeparator
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
	Director   string
	Note       string
	Size       string
}

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
	if zeros < 3 {
		return fmt.Sprintf("%s-%0*d", strings.ToUpper(id), zeros, number)
	} else {
		return fmt.Sprintf("%s-%03d", strings.ToUpper(id), number)
	}
}

func idInURL(id string, number int) string {
	return fmt.Sprintf("%s%0*d", strings.ToLower(id), zeros, number)
}

func (d *Data) scrape() error {
	if strings.Contains(d.URL, "dmm.co.jp/digital/videoc") {
		return d.dmmTypeC()
	} else if strings.Contains(d.URL, "dmm.co.jp") {
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
		if d.Date != "" && d.Director != "" {
			return
		}

		if d.Date == "" && isDateState(state) {
			d.Date = strings.ReplaceAll(strings.TrimSpace(e.Text), "/", "-")
			return
		}

		if d.Director == "" && strings.HasPrefix(state, "監督") {
			d.Director = strings.TrimSpace(e.Text)
			return
		}

		state = e.Text
	})

	c.OnHTML("tr td a", func(e *colly.HTMLElement) {
		if d.Note != "" {
			return
		}

		if strings.Contains(e.Text, "女優ベスト・総集編") || strings.Contains(e.Text, "ベスト・総集編") {
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

		d.SmallImage = pageImage(e.Attr("content"))
	})

	c.OnHTML("a[name=package-image]", func(e *colly.HTMLElement) {
		if d.LargeImage != "" {
			return
		}

		d.LargeImage = e.Attr("href")
	})

	if err := c.Visit(d.URL); err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil
		}

		return err
	}

	if d.SmallImage != "" && d.LargeImage == "" {
		d.LargeImage = strings.Replace(d.SmallImage, "ps.jpg", "pl.jpg", 1)
	}

	if d.Director == "----" {
		d.Director = ""
	}

	d.Performer = formatPerformers(performers)
	return nil
}

func (d *Data) dmmTypeC() error {
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

	state := ""
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		if d.Date != "" && d.Size != "" && d.Title != "" {
			return
		}

		text := strings.TrimSpace(e.Text)
		if d.Date == "" && strings.HasPrefix(state, "配信開始日") {
			d.Date = strings.ReplaceAll(strings.TrimSpace(text), "/", ".")
			return
		} else if d.Size == "" && strings.HasPrefix(state, "サイズ") {
			d.Size = strings.TrimSpace(e.Text)
			return
		} else if d.Title == "" && strings.HasPrefix(state, "名前") {
			d.Title = strings.TrimSpace(e.Text)
			return
		}

		state = text
	})

	c.OnHTML("meta[property=og\\:image]", func(e *colly.HTMLElement) {
		if d.SmallImage != "" {
			return
		}

		d.SmallImage = e.Attr("content")
	})

	if err := c.Visit(d.URL); err != nil {
		return err
	}

	if strings.HasSuffix(d.SmallImage, "jp.jpg") {
		d.LargeImage = d.SmallImage
		d.SmallImage = strings.Replace(d.SmallImage, "jp.jpg", "js.jpg", 1)
	}

	d.Title = fmt.Sprintf("%s~~%s", d.Title, d.Size)
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

	baseNumber, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Printf("invalid product number %s: %v\n", numberStr, err)
		return 1
	}

	if listCount != -1 {
		for i := 0; i < listCount; i++ {
			ids = append(ids, baseNumber+i)
		}
	}

	if len(ids) == 0 {
		ids = append(ids, baseNumber)
	}

	sort.Ints(ids)

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

		if err := d.scrape(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
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

		if id%10 == 0 && !noHeader {
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
