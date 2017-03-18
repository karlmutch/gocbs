package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/variadico/gocomplex/filestats"
	"github.com/variadico/gocomplex/funcstats"
	"github.com/variadico/gocomplex/pkgstats"
)

func main() {
	funcOn := flag.Bool("func", false, "Display func level stats (Default)")
	fileOn := flag.Bool("file", false, "Display file level stats")
	pkgOn := flag.Bool("pkg", false, "Display package level stats")
	flag.Parse()

	if !*funcOn && !*fileOn && !*pkgOn {
		*funcOn = true
	}

	files, err := filePaths(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	if *funcOn {
		if err := printFuncStats(files); err != nil {
			log.Fatal(err)
		}
	}

	if *fileOn {
		if err := printFileStats(files); err != nil {
			log.Fatal(err)
		}
	}

	if *pkgOn {
		if err := printPkgStats(files); err != nil {
			log.Fatal(err)
		}
	}
}

func printFuncStats(files []string) error {
	fmt.Println("statements - cyclo - nesting - params - function")

	for _, file := range files {
		stats, err := funcstats.New(file)
		if err != nil {
			return err
		}

		for _, st := range stats {
			fmt.Printf("%10d   %5d   %7d   %6d   %s:%d: %s\n",
				st.NumStmts,
				st.Complexity,
				st.MaxNest,
				st.NumParams,

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

func printPkgStats(dirs []string) error {
	fmt.Println("files - functions - types - vars - package")

	for _, dir := range dirs {
		st, err := pkgstats.New(dir)
		if err != nil {
			return err
		}

		fmt.Println(st)
	}

	return nil
}
