package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

type MyURL struct {
	ID      string
	LongURL string
}

type InStr struct {
	URL string `json:"url"`
}

type OutStr struct {
	Result string `json:"result"`
}

var myurl = []MyURL{}

var config struct {
	Host        string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"localhost:8080"`
	FileStorage string `env:"FILE_STORAGE_PATH" envDefault:"myfile"`
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The id parameter is missing", http.StatusBadRequest)
		return
	}

	fmt.Println(config)

	for i := range myurl {
		if myurl[i].ID == config.BaseURL+"/"+id {
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
	myurl = append(myurl, MyURL{config.BaseURL + "/" + id, string(url)})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(config.BaseURL + "/" + id))
}

func handlePostJSON(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var myStr InStr
	err := decoder.Decode(&myStr)
	if err != nil {
		http.Error(w, "The Body is missing", http.StatusBadRequest)
		return
	}

	fmt.Println(myStr)

	for i := range myurl {
		if myurl[i].ID == myStr.URL {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)

			subj := OutStr{myurl[i].LongURL}
			// кодируем JSON
			resp, err := json.Marshal(subj)
			if err != nil {
				http.Error(w, "The Body is missing", http.StatusBadRequest)
				return
			}
			// пишем тело ответа
			w.Write(resp)
			return
		}
	}
	http.Error(w, "The Body is missing", http.StatusBadRequest)

}

func init() {
	// Получим переменные окружения
	err := env.Parse(&config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	//	Получим данные из командной строки

	flag.StringVar(&config.Host, "a", config.Host, "host to listen on")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "baseUrl")
	flag.StringVar(&config.FileStorage, "f", config.FileStorage, "fileStorage")

	fmt.Println(config)

}

func main() {

	r := chi.NewRouter()
	r.Get("/{id}", handleGet)
	r.Post("/", handlePost)
	r.Post("/api/shorten", handlePostJSON)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(config.Host, r)
}
