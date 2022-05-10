package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
)

type typeInfo struct {
	name string
	ptr  bool
	pkg  *packages.Package
}

type handlerContext struct {
	typeInfo typeInfo
}

type requestParam struct {
	typeInfo typeInfo
}

type responseParam struct {
	typeInfo typeInfo
}

type apiHandler struct {
	pkg *ast.Package

	fnName    string
	ctx       *handlerContext
	reqParam  *requestParam
	respParam *responseParam

	err error
}

func (a *apiHandler) verify() bool {

	return true
}

type genx struct {
	path string
}

func (g *genx) resolve(dir string) {
	pkgConfig := &packages.Config{
		Mode: packages.NeedImports |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedDeps |
			packages.NeedModule |
			packages.NeedSyntax,
		Dir: dir,
	}

	pkgs, err := packages.Load(pkgConfig, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range pkgs {
		g.resolvePkg(p)
	}
}

func (g *genx) resolvePkg(pkg *packages.Package) {
	for _, f := range pkg.GoFiles {

		set := token.NewFileSet()
		astFile, err := parser.ParseFile(set, f, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		fns := g.findFunc(astFile, commentMatcher)
		for _, fn := range fns {
			fmt.Println(fn.Name)

			for _, field := range fn.Type.Params.List {
				g.resolveType(field.Type)
				fmt.Println("paramType:", field.Type)
			}

			for _, field := range fn.Type.Results.List {
				fmt.Println("resultType:", field.Type)
			}

		}
	}
}

func (g *genx) findFunc(file *ast.File, matcher FuncMatcher) []*ast.FuncDecl {

	var ret []*ast.FuncDecl

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if matcher.Match(fn) {
				ret = append(ret, fn)
				break
			}
		}
	}

	return ret
}

func (g *genx) resolveType(expr ast.Expr) *typeInfo {
	info := &typeInfo{
		name: "",
		ptr:  false,
		pkg:  nil,
	}

	return info
}

func main() {
	g := &genx{}

	g.resolve("testdata/handler")
}
