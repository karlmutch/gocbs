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

var fset = token.NewFileSet()

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

		astf, err := parser.ParseFile(fset, info.Name(), nil, 0)
		if err != nil {
			return err
		}

		funcs := getFunctions(astf)
		for _, st := range getStats(info.Name(), funcs) {
			fmt.Printf("%10d   %5d   %s\n", st.statements, st.complexity, st.function)
		}
	}

	return nil
}
