package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prog URL")
		return
	}

	url := os.Args[1]

	jar, _ := cookiejar.New(nil)

	client := http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			ForceAttemptHTTP2: true,
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bs, _ := httputil.DumpResponse(resp, true)

	fmt.Println(string(bs))
	fmt.Printf("## status=%d\n", resp.StatusCode)
	fmt.Printf("## protocol=%d.%d\n", resp.ProtoMajor, resp.ProtoMinor)
}
