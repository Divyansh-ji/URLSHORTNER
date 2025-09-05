package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type URL struct {
	ID         string    `json:"id"`
	OriginalUR string    `json:"original_uri"`
	ShortURI   string    `json:"short_uri"`
	CreatedAt  time.Time `json:"created_at"`
}

var urlDB = map[string]URL{}

func generateshorturl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher:", hasher)
	data := hasher.Sum(nil)
	fmt.Println(data)
	hash := hex.EncodeToString(data)
	fmt.Println("ENcoded string:", hash)
	fmt.Println("final string ", hash[:8])
	return hash[:8]

}
func createURL(OriginalURL string) string {
	shortURL := generateshorturl(OriginalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:         id,
		OriginalUR: OriginalURL,
		ShortURI:   shortURL,
		CreatedAt:  time.Now(),
	}
	return shortURL

}
func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}
func Rootpageurl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rootpageurl")
}
func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder((r.Body)).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func redirect(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url.OriginalUR, 302)

}

func main() {
	http.HandleFunc("/", Rootpageurl)
	http.HandleFunc("/shorturl", ShortURLHandler)
	http.HandleFunc("/redirect", redirect)
	fmt.Fprint(os.Stderr, "Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)

	}

}
