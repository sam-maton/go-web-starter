package main

import (
	"fmt"
	"log"
)

func main() {
	answers := model{}

	form := createForm(&answers)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = answers.createFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project folder '%s' created successfully!\n", answers.folder)
}
