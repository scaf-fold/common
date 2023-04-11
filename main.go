package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	path := "dddd&ddd1"
	paths := strings.Split(path, string(filepath.Separator))
	fmt.Println(paths)
}
