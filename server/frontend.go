package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type indexHandler struct {
	handler http.Handler

	filePath    string
	contentType string
}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, openErr := os.Open(h.filePath)

	if openErr != nil {
		log.Printf("Could not get index.html: %s", openErr)
		w.WriteHeader(500)
		return
	}

	statInfo, statErr := os.Stat(h.filePath)
	if statErr != nil {
		log.Printf("Could not get index.html: %s", openErr)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Length", strconv.FormatInt(statInfo.Size(), 10))
	w.Header().Add("Content-Type", h.contentType)
	w.WriteHeader(200)

	_, err := io.Copy(w, f)

	if err != nil {
		log.Printf("Could not output fronend: ", err)
	}
}
