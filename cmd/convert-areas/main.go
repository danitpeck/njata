package main

import (
	"flag"
	"fmt"
	"os"

	"njata/internal/area"
)

func main() {
	areDir := flag.String("from", "legacy/area", "Source directory with .are files")
	outputDir := flag.String("to", "areas", "Output directory for JSON files")
	flag.Parse()

	fmt.Printf("Converting .are files from %s to JSON in %s\n", *areDir, *outputDir)

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	if err := area.ConvertAreToJSON(*areDir, *outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Conversion complete!")
}
