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

	answers.extras = map[string]bool{
		"auth": false,
		"air":  false,
	}

	for _, e := range answers.selectedExtras {
		answers.extras[e] = true
	}

	err = answers.createFolders()
	if err != nil {
		log.Fatal(err)
	}

	err = answers.createBaseFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project folder '%s' created successfully!\n", answers.folder)
}
