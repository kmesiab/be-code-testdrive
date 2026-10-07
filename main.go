package main

import (
	"fmt"
	"os"

	"github.com/kmesiab/be-code-testdrive/wordfreq"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: wordfreq <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	
	freqs, err := wordfreq.WordFreq(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, freq := range freqs {
		fmt.Printf("%d %s\n", freq.Count, freq.Word)
	}
}