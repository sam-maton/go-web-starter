package main

import "fmt"

type answers struct {
	folder     string
	moduleName string
	auth       bool
	styling    bool
}

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Blue      = "\033[34m"
	Cyan      = "\033[36m"
	Magenta   = "\033[35m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgBlue    = "\033[44m"
	BgCyan    = "\033[46m"
	BgMagenta = "\033[45m"

	Bold = "\033[1m"
)

func main() {
	responses := answers{
		folder:     "my-project",
		moduleName: "github.com/you/my-project",
		auth:       true,
		styling:    true,
	}

	var (
		auth    string
		styling string
	)

	printTitle("Welcome to the Go Web Starter!")
	printLabel("Enter your project name:")

	fmt.Scan(&responses.folder)

	printLabel("Enter your Go module name:")

	fmt.Scan(&responses.moduleName)

	printLabel("Setup auth? [Y/n]")

	fmt.Scan(&auth)

	responses.auth = auth == "" || auth == "y" || auth == "Y"

	printLabel("Base styles? [Y/n]")

	fmt.Scan(&styling)

	responses.styling = styling == "" || styling == "y" || styling == "Y"

	fmt.Println(responses)
}

func printTitle(title string) {
	fmt.Println(Bold + Cyan + title + Reset)
}

func printLabel(label string) {
	fmt.Println(Magenta + "\n" + label + Reset)
}
