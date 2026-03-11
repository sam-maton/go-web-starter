package main

import (
	"os"
	"text/template"
)

func (m model) createFolders() error {
	for _, folder := range folderPaths {
		err := os.MkdirAll(m.folder+"/"+folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m model) createBaseFiles() error {
	for _, file := range baseFiles {
		filePath := m.folder + "/" + file
		err := os.WriteFile(filePath, []byte(""), 0644)
		if err != nil {
			return err
		}

		templatePath := getTemplatePath(file)
		t, err := template.ParseFiles(templatePath)
		if err != nil {
			return err
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = t.Execute(f, m.extras)
		if err != nil {
			return err
		}
	}
	return nil
}
