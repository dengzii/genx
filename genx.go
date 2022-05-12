package main

import (
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
)

func resolvePkg(pkg *packages.Package) {

	for _, f := range pkg.GoFiles {

		set := token.NewFileSet()
		astFile, err := parser.ParseFile(set, f, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		gf := NewGoFile(pkg, astFile, f)

		handlers := resolveHandlerFunc(gf, commentMatcher)
		for _, handler := range handlers {
			if len(handler.errors) > 0 {
				log.Println(handler.errors)
				continue
			}
		}
		NewGenerator(gf, handlers).Generate()

	}
}

func resolveHandlerFunc(gf *GoFile, matcher FuncMatcher) []*GoFunc {

	var handlers []*GoFunc

	gf.traversalFuncByMatch(matcher, func(handler *GoFunc) {
		handler.ResolveTypeInfo()
		if len(handler.errors) > 0 {
			log.Printf("%s: %s", handler.GetName(), handler.errors)
			return
		}
		handlers = append(handlers, handler)
	})
	return handlers
}

func main() {
	dir := "testdata/handler"

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
		resolvePkg(p)
	}
}
