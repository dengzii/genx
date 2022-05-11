package main

import (
	"go/ast"
	"strings"
)

var commentMatcher FuncMatcher = &CommentFuncMatcher{
	comment:  "//go:generate genx api",
	matchRev: true,
}

type FuncMatcher interface {
	Match(fnDecl *ast.FuncDecl) bool
}

type CommentFuncMatcher struct {
	comment  string
	matchRev bool
}

func (c *CommentFuncMatcher) Match(fnDecl *ast.FuncDecl) bool {

	if fnDecl.Doc == nil || fnDecl.Doc.List == nil || len(fnDecl.Doc.List) == 0 {
		return false
	}

	for _, comment := range fnDecl.Doc.List {
		if strings.HasPrefix(comment.Text, c.comment) {
			if fnDecl.Recv != nil && !c.matchRev {
				return false
			}
			return true
		}
	}
	return false
}
