package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func WriteJsonContent(content interface{}, w http.ResponseWriter, statusCode int) error {
	contentBytes, marshalErr := json.Marshal(content)
	if marshalErr != nil {
		return marshalErr
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(contentBytes)), 10))
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(contentBytes)
	return writeErr
}

func AssertHttpMethod(request *http.Request, w http.ResponseWriter, method string) bool {
	if request.Method != method {
		log.Printf("Got a %s request, expecting %s", request.Method, method)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(405)
		w.Write([]byte(fmt.Sprintf("Invalid method, expecting %s", method)))
		return false
	} else {
		return true
	}
}
