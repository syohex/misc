package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/atotto/clipboard"
	"github.com/playwright-community/playwright-go"
)

var ids []int
var zeros int
var withDirector bool

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

func (d *Data) scrape(browser playwright.Browser) error {
	if strings.Contains(d.URL, "dmm.co.jp") {
		return d.dmm(browser)
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

func stringPointer(s string) *string {
	return &s
}

func (d *Data) dmm(browser playwright.Browser) error {
	page, err := browser.NewPage()
	if err != nil {
		return err
	}
	page.SetDefaultTimeout(15000.0)

	page.Context().AddCookies([]playwright.OptionalCookie{
		{
			Name:   "age_check_done",
			Value:  "1",
			Path:   stringPointer("/"),
			Domain: stringPointer(".dmm.co.jp"),
		},
	})

	res, err := page.Goto(d.URL)
	if err != nil {
		return err
	}

	if res.Status() != 200 {
		return err
	}

	allPerformerButton := page.Locator("a#a_performer")
	if allPerformerButton != nil {
		count, err := allPerformerButton.Count()
		if err == nil && count != 0 {
			timeout := 5000.0 // 5 seconds
			if err := allPerformerButton.Click(playwright.LocatorClickOptions{Timeout: &timeout}); err != nil {
				fmt.Printf("failed to click all performer button: %v(%s)\n", err, d.URL)
				return err
			}

			time.Sleep(1 * time.Second)
		}
	}

	title, err := page.Locator("h1#title").TextContent()
	if err != nil {
		fmt.Printf("cannot find 'title': %v, %s(status=%d)\n", err, d.URL, res.Status())
		return err
	}

	d.Title = title

	infoEntries, err := page.Locator("tr td").All()
	if err != nil {
		fmt.Printf("cannot get product info elements: %v\n", err)
		return err
	}

	state := ""
	for _, entry := range infoEntries {
		content, err := entry.TextContent()
		if err != nil {
			fmt.Printf("cannot get text content: %v\n", err)
			return err
		}

		if d.Date != "" && d.Director != "" {
			break
		}

		if d.Date == "" && (strings.HasPrefix(state, "発売日") || strings.HasPrefix(state, "商品発売日")) {
			d.Date = strings.ReplaceAll(strings.TrimSpace(content), "/", "-")
			continue
		}

		if d.Director == "" && strings.HasPrefix(state, "監督") {
			d.Director = strings.TrimSpace(content)
			continue
		}

		state = content
	}

	genreEntries, err := page.Locator("tr td a").All()
	if err != nil {
		fmt.Printf("cannot get genre elements: %v\n", err)
		return err
	}

	for _, entry := range genreEntries {
		content, err := entry.TextContent()
		if err != nil {
			fmt.Printf("cannot get text content: %v\n", err)
			return err
		}

		if d.Note != "" {
			break
		}

		if strings.Contains(content, "女優ベスト・総集編") || strings.Contains(content, "ベスト・総集編") {
			d.Note = "総集編作品"
		}
	}

	performerEntries, err := page.Locator("td span#performer a").All()
	if err != nil {
		fmt.Printf("failed to find performer entries: %v\n", err)
		return err
	}

	var performers []string
	for _, entry := range performerEntries {
		content, err := entry.TextContent()
		if err != nil {
			fmt.Printf("cannot get text content: %v\n", err)
			return err
		}

		performers = append(performers, stripPerformer(content))
	}

	packageImageElem := page.Locator("meta[property=og\\:image]")
	smallImage, err := packageImageElem.GetAttribute("content")
	if err != nil {
		fmt.Printf("failed to find package image: %v\n", err)
		return err
	}
	d.SmallImage = pageImage(smallImage)

	largeImageElem := page.Locator("a[name=package-image]")
	largeImage, err := largeImageElem.GetAttribute("href")
	if err != nil {
		fmt.Printf("failed to find package image: %v\n", err)
		return err
	}
	d.LargeImage = largeImage

	if d.SmallImage != "" && d.LargeImage == "" {
		d.LargeImage = strings.Replace(d.SmallImage, "ps.jpg", "pl.jpg", 1)
	}

	if d.Director == "----" {
		d.Director = ""
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

	baseNumber, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Printf("invalid product number %s: %v\n", numberStr, err)
		return 1
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

	pw, err := playwright.Run()
	if err != nil {
		fmt.Printf("failed to initialize playwright: %v\n", err)
		return 1
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		fmt.Printf("failed to launch chronium")
	}
	defer browser.Close()

	baseID := idInURL(productID, baseNumber)
	baseData := &Data{
		ID:  toProductID(productID, baseNumber),
		URL: generateURL(baseURL, baseID, productID, baseNumber),
	}

	if err := baseData.scrape(browser); err != nil {
		fmt.Printf("failed to get base data %s %v\n", baseData.URL, err)
		return 1
	}

	if baseData.LargeImage == "" {
		fmt.Printf("failed to parse base data %s, data=%+v\n", baseData.URL, baseData)
		return 1
	}

	var b bytes.Buffer
	for _, id := range ids {
		d := &Data{
			ID:  toProductID(productID, id),
			URL: generateURL(baseURL, baseID, productID, id),
		}

		if err := d.scrape(browser); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scape: %v\n", err)
			break
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
