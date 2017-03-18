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

		funcs := getFunctions(f)
		for _, st := range getStats(file, funcs) {
			fmt.Printf("%10d   %5d   %7d   %s:%d: %s\n",
				st.statements,
				st.complexity,
				st.nest,

				st.file,
				fset.Position(st.pos).Line,
				st.function,
			)
		}
	}

	return nil
}
