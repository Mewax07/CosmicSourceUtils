package main

import (
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

func downloadFile(url string) string {
	fileURL, err := url2.Parse(url)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()

	return fileName
}
