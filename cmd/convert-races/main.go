package main

import (
	"flag"
	"fmt"
	"njata/internal/races"
	"os"
)

func main() {
	fromDir := flag.String("from", "legacy/races", "Source directory for .race files")
	toDir := flag.String("to", "races", "Output directory for JSON files")
	flag.Parse()

	if err := os.MkdirAll(*toDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	if err := races.ConvertRacesToJSON(*fromDir, *toDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting races: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Race conversion complete!")
}
