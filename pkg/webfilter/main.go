package webfilter

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Filter   func(io.Reader, io.Writer) error
	Title    string
	Message  string
	Filename string
}

func (h *Handler) doGet(w http.ResponseWriter, req *http.Request) error {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")

	body := fmt.Sprintf(`<html>
<head><title>%[1]s</title></head>
<body>
<h1>%[1]s</h1>
<p>%[2]s</p>
<form name="form1"
  action="/"
  enctype="multipart/form-data"
  method="post"
  accept-charset="UTF-8" >
<input type="file" name="thefile" size="48" />
<input type="submit" value="upload" />
</form>
</body>
</html>`, h.Title, h.Message)
	headers.Add("Content-Length", strconv.FormatInt(int64(len(body)), 10))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, body)
	return nil
}

func (h *Handler) doPost(w http.ResponseWriter, req *http.Request) error {

	body, _, err := req.FormFile("thefile")
	if err != nil {
		return err
	}
	defer body.Close()

	tmpFile, err := ioutil.TempFile("", h.Title)
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	h.Filter(body, tmpFile)

	tmpFile.Seek(0, os.SEEK_SET)

	headers := w.Header()
	headers.Add("Content-Disposition", "attachment; filename="+h.Filename)

	http.ServeContent(w, req, h.Filename, time.Now(), tmpFile)
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL.Path)

	var err error
	if strings.EqualFold(req.Method, "get") {
		err = h.doGet(w, req)
	} else {
		err = h.doPost(w, req)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
	}
}
