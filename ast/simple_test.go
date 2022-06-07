package ast_test

import (
	l "github.com/g-knife/common/ast"
	"go/ast"
	"go/token"
	"testing"
)

func TestSimpleFile_LoadFile(t *testing.T) {
	path := "/Users/dunbar/workspace/go_workspace/src/github.com/g-knife/common/ast/simple.go"
	sp := l.NewSimpleFile(path)
	f := sp.LoadFile()
	t.Log(ast.Print(token.NewFileSet(), f))
}
