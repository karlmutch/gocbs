package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}

func mainE() error {
	var files []string
	var err error

	if len(os.Args) == 2 {
		if os.Args[1] == "./..." {
			files, err = recursive(".")
		} else if filepath.Ext(os.Args[1]) == ".go" {
			files = append(files, os.Args[1])
		} else if strings.HasSuffix(os.Args[1], "/...") {
			p := os.Args[1]
			p = p[:len(p)-len("/...")]
			src := os.ExpandEnv("$GOPATH/src")

			files, err = recursive(filepath.Join(src, p))
		} else {
			files, err = packageFiles(os.Args[1])
		}
	} else {
		files, err = readdir(".")
	}

	if err != nil {
		return err
	}

	fmt.Println("statements - cyclo - nesting - function")
	for _, file := range files {
		if filepath.Ext(file) != ".go" {
			continue
		}

		var fset = token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return err
		}

		// foo.go:14: unreachable code
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

func recursive(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if fi == nil {
			return nil
		}
		if fi.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if filepath.Ext(p) != ".go" {
			return nil
		}

		files = append(files, p)
		return nil
	})

	return files, err
}

func packageFiles(pkg string) ([]string, error) {
	pkgPath := filepath.Join(os.ExpandEnv("$GOPATH/src/"), pkg)
	return readdir(pkgPath)
}

func readdir(p string) ([]string, error) {
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if filepath.Ext(fi.Name()) != ".go" {
			continue
		}

		files = append(files, filepath.Join(p, fi.Name()))
	}

	return files, nil
}
