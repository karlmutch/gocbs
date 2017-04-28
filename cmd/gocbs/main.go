package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/kisielk/gotool"
	"github.com/variadico/gocbs/fnstats"
	"github.com/variadico/gocbs/pkgstats"
)

var colTmpl = map[string]string{
	"name":  "{{.Name}}",
	"file":  "{{.File}}",
	"files": "{{.Files}}",
	"line":  "{{.Line}}",

	"expconst": "{{.Exported.Const}}",
	"expvar":   "{{.Exported.Var}}",
	"exptype":  "{{.Exported.Type}}",
	"expfunc":  "{{.Exported.Func}}",
	"const":    "{{.NotExported.Const}}",
	"var":      "{{.NotExported.Var}}",
	"type":     "{{.NotExported.Type}}",
	"func":     "{{.NotExported.Func}}",

	"params":  "{{.Params}}",
	"stmts":   "{{.Stmts}}",
	"cyclo":   "{{.Cyclo}}",
	"nesting": "{{.Nest}}",
}

func main() {
	pkgOn := flag.Bool("pkg", false, "Display package level stats")
	funcOn := flag.Bool("func", false, "Display function level stats (Default)")
	col := flag.String("col", "", "Column to display")
	flag.Usage = func() {
		fmt.Println(`usage: gocbs [packages]
  -func
    	Display function level stats (Default)
  -pkg
    	Display package level stats
  -o
    	Columns to display
`)
	}
	flag.Parse()

	if !*pkgOn && !*funcOn {
		*funcOn = true
	}

	paths := gotool.ImportPaths(flag.Args())
	cols := strings.Split(*col, ",")
	format := getFormat(cols)

	if *pkgOn {
		if err := printPkgStats(format, paths); err != nil {
			log.Fatal(err)
		}
	}

	if *funcOn {
		if err := printFnStats(format, paths); err != nil {
			log.Fatal(err)
		}
	}
}

func getFormat(cols []string) string {
	var buf bytes.Buffer

	for _, c := range cols {
		buf.WriteString(colTmpl[c] + " ")
	}

	return buf.String()
}

func printPkgStats(format string, paths []string) error {
	fmt.Println(pkgstats.Header)

	tmpl, err := template.New("").Parse(format)
	if err != nil {
		return err
	}

	for _, p := range paths {
		stats, err := pkgstats.New(p)
		if err != nil {
			return err
		}

		err = tmpl.Execute(os.Stdout, stats)
		if err != nil {
			return err
		}
	}

	return nil
}

func printFnStats(format string, paths []string) error {
	fmt.Println(fnstats.Header)

	tmpl, err := template.New("").Parse(format)
	if err != nil {
		return err
	}
	fmt.Println(format)

	for _, p := range paths {
		stats, err := fnstats.New(p)
		if err != nil {
			return err
		}
		for _, s := range stats {
			err = tmpl.Execute(os.Stdout, s)
			if err != nil {
				return err
			}
			fmt.Println()
		}
	}

	return nil
}
