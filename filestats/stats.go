package filestats

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Data struct {
	Name string

	Funcs int
	Types int
}

func New(filename string) (Data, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return Data{}, err
	}

	var name string
	prefix := os.ExpandEnv("$GOPATH/src/")
	if strings.HasPrefix(filename, prefix) {
		name = filename[len(prefix):]
	} else {
		name = filename
	}

	return Data{
		Name:  name,
		Funcs: numFuncs(f),
		Types: numTypes(f),
	}, nil
}

func numFuncs(f *ast.File) int {
	var points int

	ast.Inspect(f, func(n ast.Node) bool {
		if _, is := n.(*ast.FuncDecl); is {
			points++
		}

		return true
	})

	return points
}

func numTypes(f *ast.File) int {
	var points int

	ast.Inspect(f, func(n ast.Node) bool {
		if _, is := n.(*ast.TypeSpec); is {
			points++
		}

		return true
	})

	return points
}
