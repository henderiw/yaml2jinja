package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const helpMsg = `yaml2jinja converts YAML specs to Jinja templates

Usage:
    yaml2jinja < /path/to/yamlspec.yaml

Examples:
    yaml2jinja < test/example1.yaml
    yaml2jinja < test/example1.yaml > example1.go
`

func printHelp(f string) {
	helpArgs := []string{"-h", "--help", "help"}
	for _, m := range helpArgs {
		if f == m {
			fmt.Println(helpMsg)
			os.Exit(0)
		}
	}
}

func main() {
	// Read args
	if len(os.Args) > 1 {
		printHelp(os.Args[1])
	}

	// Read input from the console
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data += scanner.Text() + "\n"
	}
	//fmt.Printf(data)
	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading input:", err)
	}
	// Create yaml2jinja object and invoke Convert()
	var y2j = new(Yaml2Jinja)

	result, err := y2j.Convert([]byte(data))
	if err != nil {
		log.Fatal("Invalid YAML")
	}

	fmt.Println(result)
	return

}
