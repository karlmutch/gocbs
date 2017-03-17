package main

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
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

	for _, fn := range fns {
		s := stat{file: filepath.Base(filename)}

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

func getFunctions(f *ast.File) []*ast.FuncDecl {
	var fns []*ast.FuncDecl

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			fns = append(fns, d)
		}
	}

	return fns
}

func numStatements(fn *ast.FuncDecl) int {
	var points int

	ast.Inspect(fn, func(n ast.Node) bool {
		if _, is := n.(*ast.BlockStmt); is {
			return true
		}

		if _, is := n.(ast.Stmt); is {
			points++
		}

		return true
	})

	return points
}

func functionComplexity(fn *ast.FuncDecl) int {
	points := 1

	ast.Inspect(fn, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.IfStmt, *ast.ForStmt:
			points++
		case *ast.CaseClause:
			if t.List != nil {
				points++
			}
		case *ast.CommClause:
			if t.Comm != nil {
				points++
			}
		case *ast.BinaryExpr:
			if t.Op == token.LAND || t.Op == token.LOR {
				points++
			}
		}
		return true
	})

	return points
}

func maxNest(fn *ast.FuncDecl) int {
	return maxDepth(fn.Body) + 1
}

func maxDepth(block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}

	var depths []int

	for _, st := range block.List {
		switch t := st.(type) {
		case *ast.ForStmt:
			depths = append(depths, maxDepth(t.Body))
		case *ast.IfStmt:
			depths = append(depths, maxDepth(t.Body))
		case *ast.RangeStmt:
			depths = append(depths, maxDepth(t.Body))
		case *ast.SelectStmt:
			depths = append(depths, maxDepth(t.Body))
		case *ast.SwitchStmt:
			depths = append(depths, maxDepth(t.Body))
		case *ast.TypeSwitchStmt:
			depths = append(depths, maxDepth(t.Body))
		}
	}

	if len(depths) == 0 {
		return 0
	}

	sort.Ints(depths)
	return depths[0] + 1
}
