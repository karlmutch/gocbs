package pkgstats

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

const (
	// Header is the package stats table header.
	Header = "const - var - type - func - files - package"
)

// Tokens are the global keywords in a package.
type Tokens struct {
	Const int
	Var   int
	Type  int
	Func  int
}

// Info contains stats about a package.
type Info struct {
	// Name is the name of the package.
	Name string
	// Files is the number of files in the package.
	Files int

	// Exported contains a count of the exported tokens.
	Exported Tokens
	// NotExported contains a count of the not exported tokens.
	NotExported Tokens
}

func (inf Info) String() string {
	// const - var - type - func - files - package
	return fmt.Sprintf(
		"%5d   %3d   %4d   %4d   %5d   %s",
		inf.Exported.Const+inf.NotExported.Const,
		inf.Exported.Var+inf.NotExported.Var,
		inf.Exported.Type+inf.NotExported.Type,
		inf.Exported.Func+inf.NotExported.Func,
		inf.Files,
		inf.Name,
	)
}

// New returns info about a package.
func New(importPath string) (Info, error) {
	files, err := packageFiles(importPath)
	if err != nil {
		return Info{}, err
	}

	inf := Info{Name: importPath, Files: len(files)}
	for _, name := range files {
		if err := addFileStats(&inf, name, nil); err != nil {
			return Info{}, fmt.Errorf("failed to parse file %s: %s", name, err)
		}
	}

	return inf, nil
}

func addFileStats(inf *Info, filename string, src interface{}) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return err
	}

	for _, decl := range f.Decls {
		switch t := decl.(type) {
		case *ast.GenDecl:
			countSpecs(inf, t.Tok, t.Specs)
		case *ast.FuncDecl:
			countFuncs(inf, t)
		}
	}

	return nil
}

func countSpecs(inf *Info, tok token.Token, specs []ast.Spec) {
	for _, s := range specs {
		switch t := s.(type) {
		case *ast.ValueSpec:
			countIdents(inf, tok, t.Names)
		case *ast.TypeSpec:
			countIdents(inf, tok, []*ast.Ident{t.Name})
		}
	}
}

func countIdents(inf *Info, tok token.Token, idents []*ast.Ident) {
	for _, ident := range idents {
		if ident == nil {
			continue
		}

		exported := ast.IsExported(ident.Name)

		switch tok {
		case token.TYPE:
			if exported {
				inf.Exported.Type++
			} else {
				inf.NotExported.Type++
			}
		case token.CONST:
			if exported {
				inf.Exported.Const++
			} else {
				inf.NotExported.Const++
			}
		case token.VAR:
			if exported {
				inf.Exported.Var++
			} else {
				inf.NotExported.Var++
			}
		}
	}
}

func countFuncs(inf *Info, decl *ast.FuncDecl) {
	if decl == nil || decl.Name == nil {
		return
	}

	if ident := decl.Name; ast.IsExported(ident.Name) {
		inf.Exported.Func++
		return
	}
	inf.NotExported.Func++
}
