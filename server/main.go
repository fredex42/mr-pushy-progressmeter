package main

import (
	"log"
	"net/http"
)

type MyHttpApp struct {
	index    indexHandler
	jsbundle indexHandler
}

func main() {
	var app MyHttpApp

	app.index.filePath = "../public/index.html"
	app.index.contentType = "text/html"
	app.jsbundle.filePath = "../public/js/bundle.js"
	app.jsbundle.contentType = "application/javascript"

	http.Handle("/", app.index)
	http.Handle("/static/js/bundle.js", app.jsbundle)

	startServerErr := http.ListenAndServe(":9000", nil)

	if startServerErr != nil {
		log.Fatal(startServerErr)
	}
}
