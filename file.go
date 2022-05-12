package main

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
)

type FuncMatcher func(*ast.FuncDecl) bool

var commentMatcher FuncMatcher = func(fnDecl *ast.FuncDecl) bool {
	if fnDecl.Doc == nil || fnDecl.Doc.List == nil || len(fnDecl.Doc.List) == 0 {
		return false
	}
	for _, comment := range fnDecl.Doc.List {
		if strings.HasPrefix(comment.Text, "//go:generate genx handler") {
			if fnDecl.Recv != nil {
				return false
			}
			return true
		}
	}
	return false
}

// GoFile represents a Go source file.
type GoFile struct {
	Pkg     *packages.Package
	Imports map[string]string
	Path    string
	Name    string
	AstFile *ast.File
}

func NewGoFile(pkg *packages.Package, astFile *ast.File, path string) *GoFile {
	file := GoFile{
		Pkg:     pkg,
		AstFile: astFile,
		Path:    path,
	}
	left := strings.LastIndex(path, string(os.PathSeparator))
	if left != -1 {
		file.Name = path[left+1:]
	}
	file.Name = file.Name[:len(file.Name)-len(".go")]

	file.parseImportsName()
	return &file
}

func (g *GoFile) GetDir() string {
	index := strings.LastIndex(g.Path, string(os.PathSeparator))
	if index == -1 {
		return ""
	}
	return g.Path[:index]
}

func (g *GoFile) GetImportPkgByName(pkgName string) *packages.Package {
	pkgPath := g.Imports[pkgName]
	return g.Pkg.Imports[pkgPath]
}

func (g *GoFile) traversalFuncByMatch(matcher FuncMatcher, fn func(f *GoFunc)) {
	for _, decl := range g.AstFile.Decls {
		declFn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !matcher(declFn) {
			continue
		}
		handler := NewGoFunc(g, declFn)
		fn(handler)
	}
}

func (g *GoFile) parseImportsName() {
	g.Imports = map[string]string{}

	for _, imp := range g.AstFile.Imports {

		pkgPath := imp.Path.Value
		var pkgName string

		if imp.Name != nil && imp.Name.Name != "" {
			pkgName = imp.Name.Name
		} else {
			l := strings.LastIndex(imp.Path.Value, "/")
			pkgName = imp.Path.Value[l+1:]
			pkgName = strings.Trim(pkgName, "\"")
		}

		g.Imports[pkgName] = pkgPath
	}
}
