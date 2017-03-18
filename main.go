package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/variadico/gocomplex/filestats"
	"github.com/variadico/gocomplex/funcstats"
)

func main() {
	flag.Parse()

	files, err := filePaths(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	if err := printFuncStats(files); err != nil {
		log.Fatal(err)
	}
	if err := printFileStats(files); err != nil {
		log.Fatal(err)
	}
}

func printFuncStats(files []string) error {
	fmt.Println("statements - cyclo - nesting - function")

	for _, file := range files {
		stats, err := funcstats.New(file)
		if err != nil {
			return err
		}

		for _, st := range stats {
			fmt.Printf("%10d   %5d   %7d   %s:%d: %s\n",
				st.NumStmts,
				st.Complexity,
				st.MaxNest,

				st.File,
				st.Line,
				st.Name,
			)
		}
	}

	return nil
}

func printFileStats(files []string) error {
	fmt.Println("functions - types - file")

	for _, file := range files {
		st, err := filestats.New(file)
		if err != nil {
			return err
		}

		fmt.Printf("%9d   %5d   %s\n",
			st.Funcs,
			st.Types,
			st.Name,
		)
	}

	return nil
}
