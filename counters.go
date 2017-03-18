package main

import (
	"go/ast"
	"go/token"
	"sort"
)

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

func numStmts(fn *ast.FuncDecl) int {
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

func complexity(fn *ast.FuncDecl) int {
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

	depths := []int{0}

	for _, st := range block.List {
		switch t := st.(type) {
		case *ast.ForStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		case *ast.IfStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		case *ast.RangeStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		case *ast.SelectStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		case *ast.SwitchStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		case *ast.TypeSwitchStmt:
			depths = append(depths, maxDepth(t.Body)+1)
		}
	}

	sort.Ints(depths)
	return depths[len(depths)-1]
}
