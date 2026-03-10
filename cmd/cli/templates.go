package main

import "os"

func (m model) createFiles() error {
	for _, folder := range folderPaths {
		err := os.MkdirAll(m.folder+"/"+folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
