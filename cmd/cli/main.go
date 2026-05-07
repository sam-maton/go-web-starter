package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	answers := model{
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

	err = createProjectFiles(cwd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project files created successfully in '%s'!\n", cwd)
}
