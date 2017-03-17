package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}

func mainE() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	infos, err := ioutil.ReadDir(pwd)
	if err != nil {
		return err
	}

	fmt.Println("statements - cyclo - function")
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if filepath.Ext(info.Name()) != ".go" {
			continue
		}

		var fset = token.NewFileSet()
		f, err := parser.ParseFile(fset, info.Name(), nil, 0)
		if err != nil {
			return err
		}

		// foo.go:14: unreachable code
		funcs := getFunctions(f)
		for _, st := range getStats(info.Name(), funcs) {
			fmt.Printf("%10d   %5d   %s:%d: %s\n",
				st.statements,
				st.complexity,
				st.file,
				fset.Position(st.pos).Line,
				st.function,
			)
		}
	}

	return nil
}
