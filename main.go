package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	flag.Parse()

	files, err := filePaths(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	if err := printResults(files); err != nil {
		log.Fatal(err)
	}
}

func printResults(files []string) error {
	fmt.Println("statements - cyclo - nesting - function")

	for _, file := range files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return err
		}

		for _, st := range getFuncStats(file, getFunctions(f)) {
			fmt.Printf("%10d   %5d   %7d   %s:%d: %s\n",
				st.numStmts,
				st.complexity,
				st.maxNest,

				st.file,
				fset.Position(st.pos).Line,
				st.name,
			)
		}
	}

	return nil
}
