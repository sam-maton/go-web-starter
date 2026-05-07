package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:scaffold
var scaffoldFS embed.FS

func createProjectFiles(destination string) error {
	return fs.WalkDir(scaffoldFS, "scaffold", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel("scaffold", path)
		if err != nil {
			return err
		}

		if relativePath == "." {
			return nil
		}

		targetPath := filepath.Join(destination, relativePath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		contents, err := scaffoldFS.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, contents, 0644)
	})
}
