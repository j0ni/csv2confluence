package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Flags struct {
	inFile  string
	outFile string
	header  bool
}

func parseCli() *Flags {
	flags := new(Flags)
	flag.StringVar(&flags.inFile, "in", "table.csv", "Name of file containing CSV")
	flag.StringVar(&flags.outFile, "out", "", "Name of file to export wiki markup to")
	flag.BoolVar(&flags.header, "header", false, "Use the first row as a header")
	flag.Parse()
	return flags
}

func main() {
	// get the filenames
	flags := parseCli()

	// open the CSV file
	fi, err := os.Open(flags.inFile)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var fo *os.File
	if flags.outFile == "" {
		fo = os.Stdout
	} else {
		// open the output file
		var err error
		fo, err = os.Create(flags.outFile)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// set up CSV parser
	csvReader := csv.NewReader(fi)
	csvReader.TrailingComma = true

	// read and write
	first := true
	for {
		var separator string
		if flags.header && first {
			separator = "||"
			first = false
		} else {
			separator = "|"
		}

		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		output := fmt.Sprintf("%s%s%s\n", separator, strings.Join(fields, separator), separator)
		if _, err := fo.WriteString(output); err != nil {
			panic(err)
		}
	}
}
