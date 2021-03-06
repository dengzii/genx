package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type GoField struct {
	typeName string
	name     string

	isMap       bool
	isSlice     bool
	isInterface bool

	ptr     bool
	pkgName string
	pkgPath string

	err error
}

func NewGoType(expr ast.Expr) *GoField {
	return resolveType(expr, &GoField{})
}

func (t *GoField) resolvePkgInfo(imports map[string]string, defaultPkg string) {

	if pkgPath, ok := imports[t.pkgName]; ok {
		t.pkgPath = strings.Trim(pkgPath, "\"")
	} else {
		if t.pkgName == "" {
			t.pkgName = defaultPkg
		} else {
			t.err = fmt.Errorf("package %s not found", t.pkgName)
		}
	}
}

func (t *GoField) String() string {
	return fmt.Sprintf("{%s %s}", t.pkgName, t.typeName)
}

func resolveType(expr ast.Expr, info *GoField) *GoField {

	switch expr.(type) {

	case *ast.ArrayType:
		exp := expr.(*ast.ArrayType)
		info.isSlice = true
		resolveType(exp.Elt, info)

	case *ast.StarExpr:
		info.ptr = true
		exp := expr.(*ast.StarExpr)
		resolveType(exp.X, info)

	case *ast.SelectorExpr:
		exp := expr.(*ast.SelectorExpr)
		info.pkgName = exp.X.(*ast.Ident).Name
		resolveType(exp.Sel, info)

	case *ast.Ident:
		exp := expr.(*ast.Ident)
		info.typeName = exp.Name

	case *ast.StructType:
		fmt.Println("ArrayType")

	case *ast.MapType:
		exp := expr.(*ast.MapType)
		keyType := exp.Key.(*ast.Ident).Name
		valueType := exp.Value.(*ast.Ident).Name
		if keyType != "string" || valueType != "string" {
			info.err = fmt.Errorf("map type must be string")
		}
		info.isMap = true

	case *ast.InterfaceType:
		info.typeName = "interface{}"

	default:
		info.err = fmt.Errorf("unknown type %T", expr)
	}

	return info
}
