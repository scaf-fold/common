package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetCurGenerateFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	file, err := NewCmd().ShellExec("echo $GOFILE")
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(fmt.Sprintf("%s%c%s", dir, filepath.Separator, file), "\n", ""), nil
}

func GetGoPath() (string, error) {
	result, err := NewCmd().ShellExec("echo $GOPATH")
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(result), "\n", ""), nil
}

func GoImportFilePath(workspaceImport string) (string, error) {
	workspace, err := GetGoPath()
	if err != nil {
		return workspace, err
	}
	return fmt.Sprintf("%s%c%s%c%s", workspace, filepath.Separator, "src", filepath.Separator, workspaceImport), nil
}
