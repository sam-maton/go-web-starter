package main

import (
	"path"
	"strings"
)

func getTemplatePath(file string) string {
	ext := path.Ext(file)
	return "templates/" + strings.TrimSuffix(file, ext) + ".txt"
}
