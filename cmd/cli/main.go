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
		folder:  "my-project",
		confirm: true,
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

	err = createProjectFiles(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project files created successfully in '%s'!\n", projectPath)
}
