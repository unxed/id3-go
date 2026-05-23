package main

import (
	"fmt"
	"os"

	"github.com/unxed/id3-go"
)

func main() {
	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}

	mode := os.Args[1]
	source := os.Args[2]
	dest := os.Args[3]

	var err error
	switch mode {
	case "tojson":
		err = id3.ConvertToJSON(source, dest)
	case "toid3":
		err = id3.ConvertToID3(source, dest)
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  id3json tojson <mp3-file> <json-file>  - Convert MP3 tags to JSON technical dump")
	fmt.Println("  id3json toid3 <json-file> <mp3-file>  - Convert JSON technical dump back to MP3 tags")
}