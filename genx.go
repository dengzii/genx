package main

import (
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
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
		handlers := resolveFile(f, pkg)
		if len(handlers) == 0 {
			continue
		}
		NewGenerator(handlers).Generate()
	}
}

func resolveFile(path string, pkg *packages.Package) []*ApiHandler {

	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	gf := NewGoFile(pkg, astFile, path)
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

	return handlers
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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

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
