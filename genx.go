package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/fs"
	"log"
)

type AstFile struct {
	pkg      *ast.Package
	file     *ast.File
	typeName string
}

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

func main() {
	pkgName := ""
	fnName := ""
	pkgConfig := &packages.Config{
		Mode: packages.NeedImports |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedDeps |
			packages.NeedSyntax,
		Dir: "testdata",
	}

	pkgs, err := packages.Load(pkgConfig, "")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg.PkgPath)
		for _, f := range pkg.GoFiles {
			fmt.Println(f)
			set := token.NewFileSet()
			_, err := parser.ParseDir(set, "", func(info fs.FileInfo) bool {
				return true
			}, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	fmt.Println(pkgName, fnName)
}
