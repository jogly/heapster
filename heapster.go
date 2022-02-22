package main

import (
	_ "embed"
	"flag"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
)

//go:embed heap.gotmpl
var tmpl string

type Config struct {
	Package, Name, Type string
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := Config{}
	flag.StringVar(&cfg.Package, "pkg", "", "package name")
	flag.StringVar(&cfg.Name, "n", "", "struct name")
	flag.Parse()
	cfg.Type = flag.Args()[0]

	if cfg.Package == "" {
		cfg.Package = os.Getenv("GOPACKAGE")
	}
	if cfg.Name == "" {
		cfg.Name = extractName(cfg.Type) + "Heap"
	}
	if cfg.Type[0] != '*' && cfg.Type[0] != '[' {
		cfg.Type = "*" + cfg.Type
	}

	if err := template.Must(template.New("").Parse(tmpl)).Execute(os.Stdout, cfg); err != nil {
		panic(err)
	}
}

func extractName(s string) string {
	// Extract a name from a Go type string.
	// Example, []int -> Int
	b := strings.Builder{}

	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
		}
	}

	return strcase.ToCamel(b.String())
}
