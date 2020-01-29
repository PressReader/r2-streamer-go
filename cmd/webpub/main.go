package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/readium/r2-streamer-go/parser"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "usage: webpub PATH")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	publication, err := parser.Parse(string(inputPath))

	if err != nil {
		panic(err)
	}

	jsonData, jsonError := json.MarshalIndent(publication, "", "  ")

	if jsonError != nil {
		panic(jsonError)
	}

	fmt.Println(string(jsonData))
}
