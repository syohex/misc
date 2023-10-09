package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

var re = regexp.MustCompile(`^([^(（ ]+)`)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Usage: %s p1 p2 p3 ...\n", os.Args[0])
		return
	}

	var ws []string
	for _, a := range args {
		a = strings.TrimSpace(a)
		m := re.FindStringSubmatch(a)
		if len(m) == 0 {
			continue
		}

		ws = append(ws, fmt.Sprintf("[[%s]]", m[1]))
	}

	output := strings.Join(ws, "／")

	if err := clipboard.WriteAll(output); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)
}
