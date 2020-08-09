package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zat-kaoru-hayama/month2days"
)

type Handler struct {
	Filter func(io.Reader, io.Writer) error
}

func (h *Handler) getRequest(w http.ResponseWriter, req *http.Request) error {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")

	body := `<html>
<body>
<h1>month2days</h1>

<p>Please upload the monthly TSV-files to download the daily TSV-files converted.</p>
<form name="form1"
  action="/"
  enctype="multipart/form-data"
  method="post"
  accept-charset="UTF-8" >
<input type="file" name="tsvfile" size="48" />
<input type="submit" value="upload" />
</form>
</body>
</html>
`
	headers.Add("content-Length", strconv.FormatInt(int64(len(body)), 10))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, body)
	return nil
}

func (h *Handler) postRequest(w http.ResponseWriter, req *http.Request) error {

	body, _, err := req.FormFile("tsvfile")
	if err != nil {
		return err
	}
	defer body.Close()

	tmpFile, err := ioutil.TempFile("", "month2days")
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	h.Filter(body, tmpFile)

	tmpFile.Seek(0, os.SEEK_SET)

	headers := w.Header()
	headers.Add("Content-Disposition", "attachment; filename=output.zip")

	http.ServeContent(w, req, "output.zip", time.Now(), tmpFile)
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL.Path)

	var err error
	if strings.EqualFold(req.Method, "get") {
		err = h.getRequest(w, req)
	} else {
		err = h.postRequest(w, req)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
	}
}

func main() {
	handler := &Handler{
		Filter: func(r io.Reader, w io.Writer) error {
			storage := month2days.New()
			storage.Add(r, os.Stderr)
			storage.DumpZip(w)
			return nil
		},
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
