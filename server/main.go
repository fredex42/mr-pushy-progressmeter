package main

import (
	"log"
	"net/http"
	)

type MyHttpApp struct {
	frontend frontendHandler
}

func main() {
	var app MyHttpApp

	http.Handle("/", app.frontend)
	startServerErr := http.ListenAndServe(":9000",nil)
	if(startServerErr!=nil){
		log.Fatal(startServerErr)
	}
}