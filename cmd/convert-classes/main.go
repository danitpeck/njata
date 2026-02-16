package main

import (
	"flag"
	"fmt"
	"njata/internal/classes"
	"os"
)

func main() {
	fromDir := flag.String("from", "legacy/classes", "Source directory for .class files")
	toDir := flag.String("to", "classes", "Output directory for JSON files")
	flag.Parse()

	if err := os.MkdirAll(*toDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	if err := classes.ConvertClassesToJSON(*fromDir, *toDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting classes: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Class conversion complete!")
}
