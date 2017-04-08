package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kisielk/gotool"
	"github.com/variadico/gocbs/stats/fnstats"
	"github.com/variadico/gocbs/stats/pkgstats"
)

func main() {
	pkgOn := flag.Bool("pkg", false, "Display package level stats")
	funcOn := flag.Bool("func", false, "Display function level stats (Default)")
	flag.Usage = func() {
		fmt.Println(`usage: gocbs [packages]
  -func
    	Display function level stats (Default)
  -pkg
    	Display package level stats
`)
	}
	flag.Parse()

	if !*pkgOn && !*funcOn {
		*funcOn = true
	}

	paths := gotool.ImportPaths(flag.Args())

	if *pkgOn {
		if err := printPkgStats(paths); err != nil {
			log.Fatal(err)
		}
	}

	if *funcOn {
		if err := printFnStats(paths); err != nil {
			log.Fatal(err)
		}
	}
}

func printPkgStats(paths []string) error {
	fmt.Println(pkgstats.Header)

	for _, p := range paths {
		stats, err := pkgstats.New(p)
		if err != nil {
			return err
		}

		fmt.Println(stats)
	}

	return nil
}

func printFnStats(paths []string) error {
	fmt.Println(fnstats.Header)

	for _, p := range paths {
		stats, err := fnstats.New(p)
		if err != nil {
			return err
		}
		for _, s := range stats {
			fmt.Println(s)
		}
	}

	return nil
}
