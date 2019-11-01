package main

import (
	"log"
	"net/http"
)

type frontendHandler struct {
	handler http.Handler
}

func (h frontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		log.Printf("Could not output fronend: ", err)
	}
}

//func frontendHandlerOld(w http.ResponseWriter, r *http.Request) {
