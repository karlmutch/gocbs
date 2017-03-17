package main

import (
	"go/ast"
	"go/token"
)

type stat struct {
	function   string
	statements int
	complexity int
}

func getStats(fns []*ast.FuncDecl) []stat {
	var stats []stat

	for _, fn := range fns {
		var s stat
		if fn.Name != nil {
			s.function = fn.Name.Name
		}
		if fn.Body == nil {
			continue
		}
		s.statements = numStatements(fn)
		s.complexity = functionComplexity(fn)

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
	v := &stmts{}
	ast.Walk(v, fn)
	return v.n
}

type stmts struct {
	n int
}

func (s *stmts) Visit(node ast.Node) ast.Visitor {
	if _, is := node.(*ast.BlockStmt); is {
		return s
	}

	if _, is := node.(ast.Stmt); is {
		s.n++
	}
	return s
}

type complexity struct {
	n int
}

func (c *complexity) Visit(node ast.Node) ast.Visitor {
	switch t := node.(type) {
	case *ast.IfStmt, *ast.ForStmt:
		c.n++
	case *ast.CaseClause:
		if t.List != nil {
			c.n++
		}
	case *ast.CommClause:
		if t.Comm != nil {
			c.n++
		}
	case *ast.BinaryExpr:
		if t.Op == token.LAND || t.Op == token.LOR {
			c.n++
		}
	}

	return c
}

func functionComplexity(fn *ast.FuncDecl) int {
	v := &complexity{1}
	ast.Walk(v, fn)
	return v.n
}
