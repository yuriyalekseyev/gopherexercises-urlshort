package main

import (
	"flag"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"net/http"
	"urlshort/handlers"
)

func main() {
	ymlFileName := flag.String("file", "src/yml.yml", "a yml file with path/url pairs")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.DefaultHandler)

	db, err := buntdb.Open(":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("/db", "https://google.com", nil)
		return err
	})
	if err != nil {
		panic(err)
	}

	mapHandlerPairs := []handlers.PathUrlPair{
		{
			Path: "/g",
			Url:  "https://google.com",
		},
		{
			Path: "/gh",
			Url:  "https://github.com",
		},
		{
			Path: "/y",
			Url:  "https://yandex.ru",
		},
	}
	mapHandler := handlers.MapHandler(mapHandlerPairs, mux)

	ymlFile, err := ioutil.ReadFile(*ymlFileName)
	if err != nil {
		panic(err)
	}

	ymlHandler, err := handlers.YAMLHandler(ymlFile, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFile, err := ioutil.ReadFile("src/json.json")
	if err != nil {
		panic(err)
	}

	jsonHandler, err := handlers.JSONHandler(jsonFile, ymlHandler)
	if err != nil {
		panic(err)
	}

	dbHandler, err := handlers.DBHandler(db, jsonHandler)
	if err != nil {
		panic(err)
	}

	_ = http.ListenAndServe(":8080", dbHandler)
}
