package main

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed all:_scaffold
var scaffoldFS embed.FS

const (
	scaffoldGoModPath = "_scaffold/go.mod.txt"
)

type scaffoldTemplateData struct {
	ModuleName string
}

func createProjectFiles(destination, moduleName string) error {
	templateData := scaffoldTemplateData{ModuleName: moduleName}

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
			contents, err = renderTemplate(contents, templateData)
			if err != nil {
				return err
			}
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, contents, 0644)
	})
}

func renderTemplate(contents []byte, data scaffoldTemplateData) ([]byte, error) {
	tmpl, err := template.New("scaffold").Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return nil, err
	}

	return rendered.Bytes(), nil
}
