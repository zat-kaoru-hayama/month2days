package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/zat-kaoru-hayama/month2days"
)

var flagOutput = flag.String("o", "output.zip", "output-zip name")

func mains(args []string) error {
	if len(args) <= 0 {
		return errors.New("no data")
	}
	zipFd, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer zipFd.Close()

	return month2days.Convert(args, zipFd)
}

func main() {
	flag.Parse()
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
