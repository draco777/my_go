package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type MyURL struct {
	ID      string `json:"id"`
	LongURL string `json:"LongURL"`
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
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage string `env:"FILE_STORAGE_PATH" envDefault:"myfile"`
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "The id parameter is missing", http.StatusBadRequest)
		return
	}

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
}

func main() {

	flag.Parse()

	// Прочитаем данные из файла
	LoadDate(config.FileStorage)

	r := chi.NewRouter()
	r.Get("/{id}", handleGet)
	r.Post("/", handlePost)
	r.Post("/api/shorten", handlePostJSON)
	r.Get("/api/shorten", handleGetJSON)

	// запуск сервера

	server := &http.Server{Addr: config.Host, Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	SaveDate(config.FileStorage)
	defer cancel()
	server.Shutdown(ctx)
}
