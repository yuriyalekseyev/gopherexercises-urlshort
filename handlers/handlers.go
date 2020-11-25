package handlers

import (
	"encoding/json"
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

	pathsToUrls := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		pathsToUrls[pair.Path] = pair.Url
	}

	return MapHandler(pathsToUrls, fallback), err
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pairs []PathUrlPair
	err := yaml.Unmarshal(yml, &pairs)
	if err != nil {
		panic(err)
	}

	pathsToUrls := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		pathsToUrls[pair.Path] = pair.Url
	}

	return MapHandler(pathsToUrls, fallback), err
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		url, ok := pathsToUrls[path]
		if ok {
			fmt.Println("Redirected to " + url + " from " + path)
			http.Redirect(writer, request, url, http.StatusFound)
		} else {
			fmt.Println(path + " not found in short urls map")
			fallback.ServeHTTP(writer, request)
		}
	}
}

func DefaultHandler(writer http.ResponseWriter, _ *http.Request) {
	fmt.Println("Default handler")
	_, _ = fmt.Fprintln(writer, "You are in root")
}
