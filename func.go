package main

import "go/ast"

type GoFunc struct {
	gf      *GoFile
	astFunc *ast.FuncDecl

	receiver *GoField
	param    []*GoField
	results  []*GoField

	errors []error
}

func NewGoFunc(gf *GoFile, funcDecl *ast.FuncDecl) *GoFunc {
	return &GoFunc{
		gf:      gf,
		astFunc: funcDecl,
	}
}

func (g *GoFunc) verify() bool {
	return true
}

func (g *GoFunc) Name() string {
	return g.astFunc.Name.Name
}

func (g *GoFunc) ResolveTypeInfo() {
	g.getReceiver()
	g.getResultList()
	g.getParamList()
}

func (g *GoFunc) getParamList() []*GoField {
	if len(g.param) != 0 {
		return g.param
	}
	if g.astFunc.Type.Params != nil {
		g.param = g.getTypeInfo(g.astFunc.Type.Params)
	}
	return g.param
}

func (g *GoFunc) getResultList() []*GoField {
	if len(g.results) != 0 {
		return g.results
	}
	if g.astFunc.Type.Results != nil {
		g.results = g.getTypeInfo(g.astFunc.Type.Results)
	}
	return g.results
}

func (g *GoFunc) getReceiver() *GoField {
	if g.receiver != nil {
		return g.receiver
	}
	if g.astFunc.Recv != nil {
		r := g.getTypeInfo(g.astFunc.Recv)
		g.receiver = r[0]
	}
	return g.receiver
}

func (g *GoFunc) getTypeInfo(fields *ast.FieldList) []*GoField {

	var ret []*GoField

	for _, field := range fields.List {

		goType := NewGoType(field.Type)
		goType.resolvePkgInfo(g.gf.Imports, g.gf.Pkg.Name)

		if goType.err != nil {
			g.errors = append(g.errors, goType.err)
			break
		}

		ret = append(ret, goType)
	}
	return ret
}

func (g *GoFunc) getDepImports(pkgs map[string]string) {

	for _, goType := range g.getParamList() {
		if goType.pkgPath == "" {
			continue
		}
		pkgs[goType.pkgPath] = goType.pkgName
	}

	for _, goType := range g.getResultList() {
		if goType.pkgPath == "" {
			continue
		}
		pkgs[goType.pkgPath] = goType.pkgName
	}
}
