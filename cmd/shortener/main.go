package main

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

type MyUrl struct {
	ID      string
	LongUrl string
}

var myurl = []MyUrl{}

func handleGet(w http.ResponseWriter, r *http.Request) {
	q := strings.Replace(r.URL.Path, "/", "", -1)
	if q == "" {
		http.Error(w, "The id parameter is missing", http.StatusBadRequest)
		return
	}
	println("Получен запрос GET id = ", q)

	for i := range myurl {
		if myurl[i].ID == q {
			http.Redirect(w, r, myurl[i].LongUrl, 301)
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
	println("Получен запрос POST", string(url))

	id := strconv.Itoa(len(myurl))
	myurl = append(myurl, MyUrl{id, string(url)})

	println("Добавлен новое значние", id)

	w.Write([]byte(id))
}

// HelloWorld — обработчик запроса.
func RestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)

	}
}

func main() {

	// маршрутизация запросов обработчику
	http.HandleFunc("/", RestHandler)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
