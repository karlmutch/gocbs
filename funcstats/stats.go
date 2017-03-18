package funcstats

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

type Data struct {
	File string
	Line int
	Name string

	NumStmts   int
	Complexity int
	MaxNest    int
}

func New(filename string) ([]Data, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	var stats []Data
	prefix := os.ExpandEnv("$GOPATH/src/")
	for _, fn := range getFunctions(f) {
		var s Data

		if strings.HasPrefix(filename, prefix) {
			s.File = filename[len(prefix):]
		} else {
			s.File = filename
		}

		s.Name = funcName(fn)

		if fn.Type != nil {
			s.Line = fset.Position(fn.Type.Func).Line
		}

		if fn.Body == nil {
			continue
		}

		s.NumStmts = numStmts(fn)
		s.Complexity = complexity(fn)
		s.MaxNest = maxNest(fn)

		stats = append(stats, s)
	}

	return stats, nil
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
