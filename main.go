package month2days

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Storage struct {
	data map[string][]string
}

func New() *Storage {
	return &Storage{
		data: map[string][]string{},
	}
}

func (storage *Storage) Len() int {
	return len(storage.data)
}

func (storage *Storage) Add(r io.Reader, warn io.Writer) error {
	reader := csv.NewReader(r)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1 // field number check off
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if len(record) < 2 {
			fmt.Fprintf(warn, "too few fields in line: '%s'\n",
				strings.Join(record, "\t"))
			continue
		}
		key := record[1]
		val := strings.Join(record[2:], "\t")
		storage.data[key] = append(storage.data[key], val)
	}
	return nil
}

func (storage *Storage) DumpZip(w io.Writer) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	for date, lines := range storage.data {
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
	zipWriter.Flush()
	return nil
}

func Convert(files []string, zipFd io.Writer) error {
	storage := New()
	for _, fname := range files {
		fd, err := os.Open(fname)
		if err != nil {
			return err
		}
		err = storage.Add(fd, os.Stderr)
		fd.Close()
		if err != nil {
			return fmt.Errorf("%s: %w", fname, err)
		}
	}
	if storage.Len() <= 0 {
		return errors.New("no data")
	}
	storage.DumpZip(zipFd)
	return nil
}
