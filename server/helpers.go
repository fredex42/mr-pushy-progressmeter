package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteJsonContent(content interface{}, w http.ResponseWriter) error {
	contentBytes, marshalErr := json.Marshal(content)
	if marshalErr != nil {
		return marshalErr
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(contentBytes)), 10))
	w.WriteHeader(200)
	_, writeErr := w.Write(contentBytes)
	return writeErr
}
