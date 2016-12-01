package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"unicode/utf8"
)

//
// A poor man's awk-clone, aka
// a program to create lines from CSV columns
// Format string is string, where $0...$n indicate
// CSV columns. For instance, for the following CSV file
// basedirectory,filename
// /home/user,filename
// following format string '$0/$1' gives '/home/user/filename'
// $ can be escaped with $$
//
// should probably be rolled together with other csv manipulation
// programs and commonality be refactored, but this is quick and dirty
// for now
//

func replace_with_values(vars []string, records []string) []string {

	var res []string
	for _, v := range vars {
		if v == "$$" {
			res = append(res, v)
		} else {
			i, _ := strconv.Atoi(v[1:])
			res = append(res, records[i]) // fail with out of bounds
		}
	}
	return res
}

func merge(s1, s2 []string) string {

	s := ""
	for len(s1) > 0 && len(s2) > 0 {
		if len(s1) > 0 {
			s = s + s1[0]
			s1 = s1[1:]
		}
		if len(s2) > 0 {
			s = s + s2[0]
			s2 = s2[1:]
		}
	}

	return s
}

func process(format_string string, records []string) string {
	re := regexp.MustCompile(`(\$[0-9]+|\$\$)`)
	parts := re.Split(format_string, -1)
	vars := re.FindAllString(format_string, -1)
	substituted := replace_with_values(vars, records)
	return merge(parts, substituted)
}

func main() {

	separator := flag.String("F", ";", "Field separator character")
	format_string := flag.String("s", "", "Format string")
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

	for _, record := range records {
		fmt.Fprintf(os.Stdout, "%v\n", process(*format_string, record))
	}
}
