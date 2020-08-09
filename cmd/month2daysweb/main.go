package main

import (
	"io"
	"os"

	"github.com/zat-kaoru-hayama/month2days"
	"github.com/zat-kaoru-hayama/month2days/pkg/webfilter"
)

func main() {
	h := &webfilter.Handler{
		Filter: func(r io.Reader, w io.Writer) (string, error) {
			storage := month2days.New()
			err := storage.Add(r, os.Stderr)
			if err != nil {
				return "", err
			}
			err = storage.DumpZip(w)
			if err != nil {
				return "", err
			}
			return "output.zip", nil
		},
		Title:   "month2days",
		Message: "Please upload the monthly TSV-files to download the daily TSV-files converted.",
	}
	h.Run(8000)
}
