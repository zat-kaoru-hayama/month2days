package main

import (
	"fmt"
	"os"

	"github.com/sqweek/dialog"
	"github.com/zat-kaoru-hayama/month2days"
)

func mains() error {
	inputPath, err := dialog.File().Filter("(68xxxx)YYYYMM.tsv", "tsv").Title("Please select monthly-TSV file").Load()
	if err != nil {
		return err
	}
	outputPath, err := dialog.File().Filter("output.zip", "zip").Title("Output").Save()
	if err != nil {
		return err
	}
	zipFd, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFd.Close()

	return month2days.Convert([]string{inputPath}, zipFd)
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
