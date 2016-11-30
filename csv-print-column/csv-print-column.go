package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func main() {

	separator := flag.String("F", ";", "Field separator character")
	column := flag.Int("c", 0, "Column to print, starting from zero")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] input-file\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	infile, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open input file: %v\n", err)
		os.Exit(1)
	}
	defer infile.Close()

	r := csv.NewReader(infile)
	ru, _ := utf8.DecodeRuneInString(*separator)
	r.Comma = ru

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read records: %v\n", err)
		os.Exit(1)
	}

	if len(records) < 1 {
		os.Exit(0)
	}
	if (*column > (len(records[0]) - 1)) || (*column < 0) {
		fmt.Fprintf(os.Stderr, "Column %d out of range (0-%d)\n", *column, len(records[0])-1)
		os.Exit(1)
	}

	for _, record := range records {
		fmt.Fprintf(os.Stdout, "%v\n", record[*column])
	}
}
