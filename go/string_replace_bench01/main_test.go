package main

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkRegexpReplace(b *testing.B) {
	s := "%[]^_"
	r := regexp.MustCompile(`%|\[|\]|\^|\_`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReplaceAllString(s, `\$0`)
	}
}

func BenchmarkStringReplacer(b *testing.B) {
	s := "%[]^_"
	r := strings.NewReplacer("%", `\%`, "[", `\[`, "]", `\]`, "^", `\^`, "_", `\_`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Replace(s)
	}
}

func TestRegexpReplace(t *testing.T) {
	s := "%[]^_"
	expected := `\%\[\]\^\_`
	r := regexp.MustCompile(`%|\[|\]|\^|\_`)
	got := r.ReplaceAllString(s, `\$0`)
	if got != expected {
		t.Fatalf("Got: %s, expected %s", got, expected)
	}
}

func TestStringReplacer(t *testing.T) {
	s := "%[]^_"
	expected := `\%\[\]\^\_`
	r := strings.NewReplacer("%", `\%`, "[", `\[`, "]", `\]`, "^", `\^`, "_", `\_`)
	got := r.Replace(s)
	if got != expected {
		t.Fatalf("Got: %s, expected %s", got, expected)
	}
}
