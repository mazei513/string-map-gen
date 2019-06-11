package main

import (
	"fmt"
	"os"

	"github.com/mazei513/string-map-gen/generator"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("invalid arguments given")
		fmt.Println("example usage:")
		fmt.Println("go run github.com/mazei513/string-map-gen gopher.go Gopher'")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := generator.Generate(os.Args[2], wd, os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
