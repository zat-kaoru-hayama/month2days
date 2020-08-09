package main

import (
	"archive/zip"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var flagOutput = flag.String("o", "output.zip", "output-zip name")

func mains(args []string) error {
	data := map[string][]string{}
	for _, fname := range args {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		reader := csv.NewReader(fd)
		reader.Comma = '\t'
		reader.FieldsPerRecord = -1 // field number check off
		for {
			record, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					fd.Close()
					return fmt.Errorf("%s: %w", fname, err)
				}
				break
			}
			if len(record) < 2 {
				fmt.Fprintf(os.Stderr, "%s: too few fields in line: '%s'\n",
					fname,
					strings.Join(record, "\t"))
				continue
			}
			key := record[1]
			val := strings.Join(record[2:], "\t")
			data[key] = append(data[key], val)
		}
		fd.Close()
	}
	zipFd, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer zipFd.Close()
	zipWriter := zip.NewWriter(zipFd)
	defer zipWriter.Close()
	for date, lines := range data {
		tsvName := strings.Replace(date, "/", "", -1) + ".tsv"
		fd, err := zipWriter.Create(tsvName)
		if err != nil {
			return err
		}
		fmt.Fprintf(fd, "\uFEFF; COMPUTERNAME=%s USERNAME=%s\r\n",
			os.Getenv("COMPUTERNAME"), os.Getenv("USERNAME"))

		for _, line := range lines {
			fmt.Fprintf(fd, "%s\r\n", line)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
