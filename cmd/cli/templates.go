package main

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:_scaffold
var scaffoldFS embed.FS

const (
	scaffoldGoModPath = "_scaffold/go.mod.txt"
	goModulePrefix    = "module "
	scaffoldModule    = "github.com/sam-maton/go-web-starter-baseline"
)

func createProjectFiles(destination, moduleName string) error {
	moduleName = strings.TrimSpace(moduleName)

	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	return fs.WalkDir(scaffoldFS, "_scaffold", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel("_scaffold", path)
		if err != nil {
			return err
		}

		if relativePath == "." {
			return nil
		}

		if relativePath == "go.mod.txt" {
			relativePath = "go.mod"
		}

		targetPath := filepath.Join(destination, relativePath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		contents, err := scaffoldFS.ReadFile(path)
		if err != nil {
			return err
		}

		if path == scaffoldGoModPath {
			contents = replaceGoModule(contents, moduleName)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, contents, 0644)
	})
}

func replaceGoModule(contents []byte, moduleName string) []byte {
	lines := bytes.Split(contents, []byte("\n"))

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(string(line))
		if trimmedLine == goModulePrefix+scaffoldModule {
			lines[i] = []byte(goModulePrefix + moduleName)
			return bytes.Join(lines, []byte("\n"))
		}
	}

	return contents
}
