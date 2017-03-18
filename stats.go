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

		if fn.Recv != nil {
			if len(fn.Recv.List) == 1 {
				expr := fn.Recv.List[0].Type
				if se, is := expr.(*ast.StarExpr); is {
					if v, is := se.X.(*ast.Ident); is {
						s.function = "*" + v.Name + "."
					}
				}
				if v, is := expr.(*ast.Ident); is {
					s.function = v.Name + "."
				}
			}
		}

		if fn.Name != nil {
			s.function += fn.Name.Name
		}

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
