package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func urlsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		urlsPostHandler(w, r)
		http.Error(w, "", http.StatusNotFound)
	}
}

func urlsPostHandler(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		http.Error(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
	}
	var urlsList []string
	err := json.NewDecoder(r.Body).Decode(&urlsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	response := map[string][]byte{}
	for _, url := range urlsList {
		body, err := getBody(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		response[url] = body
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return
}

func getBody(url string) ([]byte, error) {
	urlRequest, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(urlRequest.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	http.HandleFunc("/urls/", urlsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
