package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

type PathUrlPair struct {
	Path string
	Url  string
}

func JSONHandler(jsonRaw []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pairs []PathUrlPair
	err := json.Unmarshal(jsonRaw, &pairs)
	if err != nil {
		panic(err)
	}

	return MapHandler(pairs, fallback), err
}

func YAMLHandler(ymlRaw []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pairs []PathUrlPair
	err := yaml.Unmarshal(ymlRaw, &pairs)
	if err != nil {
		panic(err)
	}

	return MapHandler(pairs, fallback), err
}

func MapHandler(pairs []PathUrlPair, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		url, err := findUrl(pairs, path)
		if err != nil {
			fallback.ServeHTTP(writer, request)
		} else {
			fmt.Println("Redirected to " + url + " from " + path)
			http.Redirect(writer, request, url, http.StatusFound)
		}
	}
}

func DefaultHandler(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Default handler")
	_, _ = fmt.Fprintln(writer, "You are in root")
}

func findUrl(pairs []PathUrlPair, path string) (url string, err error) {
	for _, pair := range pairs {
		if path == pair.Path {
			return pair.Url, nil
		}
	}
	return "", errors.New("path " + path + " not found")
}
