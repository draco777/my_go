package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
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
	Host        string `env:"SERVER_ADDRESS" envDefault:"localhost"`
	Port        string `env:"PORT" envDefault:"8080"`
	BaseUrl     string `env:"BASE_URL" envDefault:"localhost"`
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

func handlePostJSON(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var my_str InStr
	err := decoder.Decode(&my_str)
	if err != nil {
		http.Error(w, "The Body is missing", http.StatusBadRequest)
		return
	}

	for i := range myurl {
		if myurl[i].ID == my_str.URL {
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

	fmt.Println(config)

	//
	flag.StringVar(&config.Port, "port", "8080", "port to listen on")
	flag.StringVar(&config.Host, "a", "localhost", "host to listen on")
	flag.StringVar(&config.BaseUrl, "b", "localhost", "baseUrl")
	flag.StringVar(&config.FileStorage, "f", "", "fileStorage")

	// Проверим что файл хранения не задан
	if config.FileStorage == "" {
		file, err := os.CreateTemp(os.TempDir(), "myfile")
		if err != nil {
			panic(err)
		}
		config.FileStorage = file.Name()

		defer file.Close()
	}

	fmt.Println(config.FileStorage)
}

func main() {

	r := chi.NewRouter()
	r.Get("/{id}", handleGet)
	r.Post("/", handlePost)
	r.Post("/api/shorten", handlePostJSON)

	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(config.Host+":"+"8080", r)
}
