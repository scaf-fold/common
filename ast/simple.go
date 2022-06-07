package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type SimpleFile struct {
	Path string
}

func NewSimpleFile(path string) *SimpleFile {
	return &SimpleFile{
		path,
	}
}

func (s *SimpleFile) LoadFile() (file *ast.File) {
	if f, err := os.Open(s.Path); err == nil {
		name := filepath.Base(s.Path)
		fset := token.NewFileSet()
		file, err = parser.ParseFile(fset, name, f, parser.ParseComments)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	return
}
