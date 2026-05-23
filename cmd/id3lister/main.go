package main

import (
	"fmt"
	"io"
	"os"

	"github.com/unxed/id3-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: id3lister <filename>")
		os.Exit(1)
	}

	if err := listFile(os.Args[1], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func listFile(filename string, w io.Writer) error {
	file, err := id3.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(w, "File: %s\n", filename)
	fmt.Fprintf(w, "Version: %s\n", file.Version())
	fmt.Fprintf(w, "Title: %s\n", file.Title())
	fmt.Fprintf(w, "Artist: %s\n", file.Artist())
	fmt.Fprintf(w, "Album: %s\n", file.Album())
	fmt.Fprintf(w, "Year: %s\n", file.Year())
	fmt.Fprintf(w, "Genre: %s\n", file.Genre())
	fmt.Fprintf(w, "Comments: %v\n", file.Comments())

	frames := file.AllFrames()
	if len(frames) > 0 {
		fmt.Fprintln(w, "Frames:")
		for _, frame := range frames {
			fmt.Fprintf(w, "  [%s] %s\n", frame.Id(), frame.String())
		}
	}
	return nil
}