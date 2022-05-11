package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
)

type apiHandler struct {
	fnName string

	receiver *typeInfo
	param    []*typeInfo
	results  []*typeInfo

	err error
}

func (a *apiHandler) verify() bool {

	return true
}

type goFile struct {
	pkg     *packages.Package
	imports map[string]string
	astFile *ast.File
}

type genx struct {
	path string

	ts *typeResolver
}

func newGenx() *genx {
	gex := &genx{
		ts: &typeResolver{},
	}
	return gex
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

		gf := &goFile{
			pkg:     pkg,
			imports: g.getImports(astFile),
			astFile: astFile,
		}

		handlers := g.resolveHandlerFunc(gf, commentMatcher)
		for _, handler := range handlers {
			fmt.Println(handler.fnName)
			if handler.err != nil {
				log.Fatal(handler.err)
			}
		}
	}
}

func (g *genx) getTypeInfo(fields *ast.FieldList, gf *goFile) ([]*typeInfo, error) {

	var ret []*typeInfo
	var err error

	for _, field := range fields.List {

		info := g.ts.resolve(field.Type)
		info.resolvePkgInfo(gf.imports, gf.pkg.Name)
		if info.err != nil {
			err = info.err
			break
		}

		ret = append(ret, info)
	}
	return ret, err
}

func (g *genx) getImports(file *ast.File) map[string]string {

	imports := map[string]string{}

	for _, imp := range file.Imports {

		pkgPath := imp.Path.Value
		var pkgName string

		if imp.Name != nil && imp.Name.Name != "" {
			pkgName = imp.Name.Name
		} else {
			l := strings.LastIndex(imp.Path.Value, "/")
			pkgName = imp.Path.Value[l+1:]
			pkgName = strings.Trim(pkgName, "\"")
		}

		imports[pkgName] = pkgPath
	}
	return imports
}

func (g *genx) resolveHandlerFunc(gf *goFile, matcher FuncMatcher) []*apiHandler {

	var handlers []*apiHandler

	for _, decl := range gf.astFile.Decls {

		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if !matcher.Match(fn) {
			continue
		}

		handler := &apiHandler{
			fnName:  fn.Name.Name,
			param:   []*typeInfo{},
			results: []*typeInfo{},
			err:     nil,
		}
		handlers = append(handlers, handler)

		if fn.Recv != nil {
			receiver, err := g.getTypeInfo(fn.Recv, gf)
			if err != nil {
				handler.err = err
				continue
			} else {
				handler.receiver = receiver[0]
			}
		}
		if fn.Type.Params != nil {
			handler.param, handler.err = g.getTypeInfo(fn.Type.Params, gf)
		}
		if handler.err != nil {
			continue
		}
		if fn.Type.Results != nil {
			handler.results, handler.err = g.getTypeInfo(fn.Type.Results, gf)
		}
	}

	return handlers
}

func main() {
	g := newGenx()

	g.resolve("testdata/handler")
}
