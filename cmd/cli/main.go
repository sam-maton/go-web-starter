package main

import (
	"fmt"
	"log"
)

func main() {
	answers := model{
		folder: "test-folder",
	}

	form := createForm(&answers)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answers)
}
