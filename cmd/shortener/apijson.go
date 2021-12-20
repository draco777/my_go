package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type APIError struct {
	Error string `json:"error"`
}

func handlePostJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var myStr InStr
	err := decoder.Decode(&myStr)
	if err != nil {
		b := APIError{"The Body is missing"}
		resp, _ := json.Marshal(b)
		w.Write(resp)
		http.Error(w, "1 The Body is missing", http.StatusBadRequest)
		return
	}

	id := strconv.Itoa(len(myurl))
	myurl = append(myurl, MyURL{config.BaseURL + "/" + id, myStr.URL})
	w.WriteHeader(http.StatusCreated)
	subj := OutStr{config.BaseURL + "/" + id}
	// кодируем JSON
	resp, _ := json.Marshal(subj)
	w.Write(resp)
}

func handleGetJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var myStr InStr
	err := decoder.Decode(&myStr)
	if err != nil {
		b := APIError{"The Body is missing"}
		resp, _ := json.Marshal(b)
		w.Write(resp)
		http.Error(w, "1 The Body is missing", http.StatusBadRequest)
		return
	}

	for i := range myurl {
		if myurl[i].ID == myStr.URL {

			subj := OutStr{myurl[i].LongURL}
			// кодируем JSON
			resp, err := json.Marshal(subj)
			if err != nil {
				b := APIError{"The Body is missing"}
				resp, _ := json.Marshal(b)
				w.Write(resp)
				http.Error(w, "2 The Body is missing", http.StatusBadRequest)
				return
			}
			// пишем тело ответа
			w.Write(resp)
			http.Redirect(w, r, myurl[i].LongURL, http.StatusTemporaryRedirect)
			return
		}
	}
	b := APIError{"The Body is missing"}
	resp, _ := json.Marshal(b)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}
