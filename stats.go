package main

import (
	"go/ast"
	"go/token"
	"os"
	"strings"
)

type stat struct {
	file     string
	pos      token.Pos
	function string

	statements int
	complexity int
	nest       int
}

func getStats(filename string, fns []*ast.FuncDecl) []stat {
	var stats []stat
	prefix := os.ExpandEnv("$GOPATH/src/")

	for _, fn := range fns {
		var s stat

		if strings.HasPrefix(filename, prefix) {
			s.file = filename[len(prefix):]
		} else {
			s.file = filename
		}

		s.function = funcName(fn)

		if fn.Type != nil {
			s.pos = fn.Type.Func
		}

		if fn.Body == nil {
			continue
		}

		s.statements = numStatements(fn)
		s.complexity = functionComplexity(fn)
		s.nest = maxNest(fn)

		stats = append(stats, s)
	}

	return stats
}

func funcName(fn *ast.FuncDecl) string {
	// Regular function.
	if fn.Recv == nil && fn.Name != nil {
		return fn.Name.Name
	}

	if len(fn.Recv.List) != 1 {
		return ""
	}

	var name string
	expr := fn.Recv.List[0].Type

	// Method on a pointer receiver.
	if se, is := expr.(*ast.StarExpr); is {
		if v, is := se.X.(*ast.Ident); is {
			name = "*" + v.Name + "."
		}
	}

	// Method on a value receiver.
	if v, is := expr.(*ast.Ident); is {
		name = v.Name + "."
	}

	if fn.Name != nil {
		name += fn.Name.Name
	}

	return name
}
