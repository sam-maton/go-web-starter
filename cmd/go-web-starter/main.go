package main

import (
"fmt"
"log"
"os"
"path/filepath"
"strings"
)

func main() {
answers := model{
folder:     "my-project",
moduleName: "github.com/you/my-project",
enableAuth: true,
confirm:    true,
}

form := createForm(&answers)

err := form.Run()
if err != nil {
log.Fatal(err)
}

if !answers.confirm {
fmt.Println("Project setup cancelled.")
return
}

cwd, err := os.Getwd()
if err != nil {
log.Fatal(err)
}

projectPath := filepath.Join(cwd, strings.TrimSpace(answers.folder))
moduleName := strings.TrimSpace(answers.moduleName)

err = createProjectFiles(projectPath, moduleName, answers.enableAuth)
if err != nil {
log.Fatal(err)
}

fmt.Printf("Project files created successfully in '%s'!\n", projectPath)
}
