package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Usage: %s p1 p2 p3 ...", os.Args[0])
		return
	}

	var ws []string
	for _, a := range args {
		ws = append(ws, fmt.Sprintf("[[%s]]", a))
	}

	output := strings.Join(ws, "／")

	if err := clipboard.WriteAll(output); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)
}
