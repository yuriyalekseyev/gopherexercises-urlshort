package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"urlshort/handlers"
)

func main() {
	ymlFileName := flag.String("file", "pairs.yml", "a yml file with path/url pairs")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.DefaultHandler)

	pathsToUrls := map[string]string{
		"/g":  "https://google.com",
		"/gh": "https://github.com",
		"/y":  "https://yandex.ru",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	ymlFile, err := ioutil.ReadFile(*ymlFileName)
	if err != nil {
		panic(err)
	}

	ymlHandler, err := handlers.YAMLHandler(ymlFile, mapHandler)
	if err != nil {
		panic(err)
	}

	_ = http.ListenAndServe(":8080", ymlHandler)
}
