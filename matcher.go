package main

import (
	"go/ast"
	"strings"
)

var commentMatcher = &CommentFuncMatcher{comment: "//genx:api"}

type FuncMatcher interface {
	Match(fnDecl *ast.FuncDecl) bool
}

type CommentFuncMatcher struct {
	comment string
}

func (c *CommentFuncMatcher) Match(fnDecl *ast.FuncDecl) bool {
	for _, comment := range fnDecl.Doc.List {
		if strings.HasPrefix(comment.Text, c.comment) {
			return true
		}
	}
	return false
}
