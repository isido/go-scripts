package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//
// A script to match filenames and item directories
// in SAF archives. Should probably be rolled into
// saf-archiver later on
//

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %v bitstream-name-file item-file bitstream-file-dir\n",
			filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\tbitstream-name-file\tfile containing relative paths to bitstreams, one per line\n")
		fmt.Fprintf(os.Stderr, "\titem-file\tfile containing relatives paths to item directories, one per line, corresponding to bitstreams\n")
		os.Exit(1)
	}

	bitstream_infile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open bitstream name file: %v\n", err)
		os.Exit(1)
	}
	defer bitstream_infile.Close()

	item_infile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open item path file: %v\n", err)
		os.Exit(1)
	}
	defer item_infile.Close()

	bitstreams := bufio.NewReader(bitstream_infile)
	items := bufio.NewReader(item_infile)

	found_bitstreams := 0
	missing_bitstreams := 0

	for {
		// ignore possible long lines
		bitstream_b, _, err := bitstreams.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read line from bitstream file: %v!\n", err)
			os.Exit(1)
		}
		bitstream := string(bitstream_b)

		item_b, _, err := items.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read line from item path file: %v\n", err)
		}
		item := string(item_b)

		if _, err := os.Stat(bitstream); err == nil {
			fmt.Fprintf(os.Stdout, "add-file -r \"%v\" %v\n", bitstream, item)
			found_bitstreams = found_bitstreams + 1
		} else {
			fmt.Fprintf(os.Stderr, "Cannot find bitstream '%v'\n", bitstream)
			missing_bitstreams = missing_bitstreams + 1
		}
	}
	fmt.Fprintf(os.Stderr, "Found items %d\n", found_bitstreams)
	fmt.Fprintf(os.Stderr, "Missing items %d\n", missing_bitstreams)
}
