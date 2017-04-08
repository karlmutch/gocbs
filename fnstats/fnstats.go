package fnstats

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

const (
	// Header is the function stats table header.
	Header = "params - stmts - cyclo - nest - func"
)

// Info contains the function stats.
type Info struct {
	// Name is the name of the function.
	Name string
	// File is the file where the function lives.
	File string
	// Line is where the function starts.
	Line int

	Params int
	Stmts  int
	Cyclo  int
	Nest   int
}

func (inf Info) String() string {
	return fmt.Sprintf(
		"%6d   %5d   %5d   %4d   %s",
		inf.Params,
		inf.Stmts,
		inf.Cyclo,
		inf.Nest,

		fmt.Sprintf("%s:%d %s", inf.File, inf.Line, inf.Name),
	)
}

// New returns a slice of information of each function in a package.
func New(importPath string) ([]Info, error) {
	files, err := packageFiles(importPath)
	if err != nil {
		return nil, err
	}

	var infs []Info
	fset := token.NewFileSet()
	for _, fname := range files {
		f, err := parser.ParseFile(fset, fname, nil, 0)
		if err != nil {
			return nil, err
		}

		for _, fn := range getFuncs(f) {
			inf := Info{
				File: fname,
				Name: funcName(fn),
			}
			if fn.Type != nil {
				inf.Line = fset.Position(fn.Type.Func).Line
			}

			countProps(&inf, fn)

			infs = append(infs, inf)
		}
	}

	return infs, nil
}

func getFuncs(f *ast.File) []*ast.FuncDecl {
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
	if fn.Recv == nil && fn.Name != nil {
		// Regular function.
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

func countProps(inf *Info, fn *ast.FuncDecl) {
	if fn.Body == nil {
		return
	}

	inf.Params = countParams(fn)
	inf.Stmts = countStmts(fn)
	inf.Cyclo = countCyclo(fn)
	inf.Nest = countNest(fn)
}

func countParams(fn *ast.FuncDecl) int {
	if fn.Type == nil {
		return 0
	}

	if fn.Type.Params == nil {
		return 0
	}

	return len(fn.Type.Params.List)
}

func countStmts(fn *ast.FuncDecl) int {
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

func countCyclo(fn *ast.FuncDecl) int {
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

func countNest(fn *ast.FuncDecl) int {
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
