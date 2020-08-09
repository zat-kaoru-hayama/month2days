package webfilter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/toqueteos/webbrowser"
)

type Handler struct {
	Filter  func(io.Reader, io.Writer) (string, error)
	Title   string
	Message string
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

	var buf bytes.Buffer
	fname, err := h.Filter(body, &buf)
	if err != nil {
		return err
	}

	headers := w.Header()
	headers.Add("Content-Disposition", "attachment; filename="+fname)

	http.ServeContent(w, req, fname, time.Now(), bytes.NewReader(buf.Bytes()))
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

func (h *Handler) Run(portNo int) {
	port := fmt.Sprintf(":%d", portNo)
	service := &http.Server{
		Addr:           port,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		webbrowser.Open("http://127.0.0.1" + port)
	}()
	service.ListenAndServe()
	service.Close()
}
