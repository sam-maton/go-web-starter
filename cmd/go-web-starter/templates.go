package main

import (
"bytes"
"embed"
"go/format"
"io/fs"
"os"
"path/filepath"
"text/template"
)

//go:embed all:_scaffold
var scaffoldFS embed.FS

// authOnlyFiles lists scaffold paths (relative to _scaffold) that should only
// be included when auth is enabled.
var authOnlyFiles = map[string]bool{
"ui/html/pages/login.html":               true,
"ui/html/pages/signup.html":              true,
"internal/models/users.go":               true,
"sql/schema/00001_create_users_table.sql": true,
}

type scaffoldTemplateData struct {
ModuleName string
EnableAuth bool
}

func createProjectFiles(destination, moduleName string, enableAuth bool) error {
templateData := scaffoldTemplateData{
ModuleName: moduleName,
EnableAuth: enableAuth,
}

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

// Skip auth-only files when auth is disabled.
if !enableAuth && authOnlyFiles[filepath.ToSlash(relativePath)] {
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

// Render every scaffold file as a template using [[ ]] delimiters so
// that Go/HTML template syntax ({{ }}) is left untouched.
contents, err = renderTemplate(contents, templateData)
if err != nil {
return err
}

// Auto-format generated Go source files.
if filepath.Ext(relativePath) == ".go" {
if formatted, fmtErr := format.Source(contents); fmtErr == nil {
contents = formatted
}
}

if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
return err
}

return os.WriteFile(targetPath, contents, 0644)
})
}

func renderTemplate(contents []byte, data scaffoldTemplateData) ([]byte, error) {
tmpl, err := template.New("scaffold").Delims("[[", "]]").Parse(string(contents))
if err != nil {
return nil, err
}

var rendered bytes.Buffer
if err := tmpl.Execute(&rendered, data); err != nil {
return nil, err
}

return rendered.Bytes(), nil
}
