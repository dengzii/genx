package main

import (
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
)

var (
	ignores = []string{"_genx.go", "_test.go", "_example.go"}
)

func resolvePkg(pkg *packages.Package) {

	for _, f := range pkg.GoFiles {

		ignore := false
		for _, s := range ignores {
			if strings.HasSuffix(f, s) {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}

		set := token.NewFileSet()
		astFile, err := parser.ParseFile(set, f, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		gf := NewGoFile(pkg, astFile, f)

		funcs := resolveHandlerFunc(gf, commentMatcher)

		handlers := make([]*ApiHandler, 0, len(funcs))

		for _, fn := range funcs {
			if len(fn.errors) > 0 {
				log.Println(fn.errors)
				continue
			}
			handler, err := NewApiHandler(fn)
			if err != nil {
				log.Println(err)
				continue
			}
			handlers = append(handlers, handler)
		}
		NewGenerator(gf, handlers).Generate()

	}
}

func resolveHandlerFunc(gf *GoFile, matcher FuncMatcher) []*GoFunc {

	var handlers []*GoFunc

	gf.traversalFuncByMatch(matcher, func(handler *GoFunc) {
		handler.ResolveTypeInfo()
		if len(handler.errors) > 0 {
			log.Printf("%s: %s", handler.Name(), handler.errors)
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
