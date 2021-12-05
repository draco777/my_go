package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

type MyURL struct {
	ID      string
	LongURL string
}

var myurl = []MyURL{}

func handleGet(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The id parameter is missing", http.StatusBadRequest)
		return
	}

	for i := range myurl {
		if myurl[i].ID == id {
			http.Redirect(w, r, myurl[i].LongURL, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Error(w, "Плохой запрос", http.StatusBadRequest)

}

func handlePost(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "The Body is missing", http.StatusBadRequest)
		return
	}

	id := strconv.Itoa(len(myurl))
	myurl = append(myurl, MyURL{id, string(url)})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + id))
}

func main() {
	r := chi.NewRouter()
	r.Get("/{id}", handleGet)
	r.Post("/", handlePost)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", r)
}
