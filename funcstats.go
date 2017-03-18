package main

import (
	"go/ast"
	"go/token"
	"os"
	"strings"
)

type funcStat struct {
	name     string
	pos      token.Pos
	funcName string

	numStmts   int
	complexity int
	maxNest    int
}

func getFuncStats(filename string, fns []*ast.FuncDecl) []funcStat {
	var stats []funcStat
	prefix := os.ExpandEnv("$GOPATH/src/")

	for _, fn := range fns {
		var s funcStat

		if strings.HasPrefix(filename, prefix) {
			s.name = filename[len(prefix):]
		} else {
			s.name = filename
		}

		s.funcName = funcName(fn)

		if fn.Type != nil {
			s.pos = fn.Type.Func
		}

		if fn.Body == nil {
			continue
		}

		s.numStmts = numStmts(fn)
		s.complexity = complexity(fn)
		s.maxNest = maxNest(fn)

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
