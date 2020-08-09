package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/zat-kaoru-hayama/month2days"
	"github.com/zat-kaoru-hayama/month2days/pkg/webfilter"
)

func main() {
	handler := &webfilter.Handler{
		Filter: func(r io.Reader, w io.Writer) error {
			storage := month2days.New()
			storage.Add(r, os.Stderr)
			storage.DumpZip(w)
			return nil
		},
		Title:    "month2days",
		Message:  "Please upload the monthly TSV-files to download the daily TSV-files converted.",
		Filename: "output.zip",
	}

	service := &http.Server{
		Addr:           ":8000",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	service.ListenAndServe()
	service.Close()
}
