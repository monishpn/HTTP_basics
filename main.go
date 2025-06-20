package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func hellohandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

}

type Record struct {
	Id        int
	Task      string
	Completed bool
}

var m = make(map[int][]byte)
var i int = 0

func addTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	rec := Record{i, string(msg), false}

	json_rec, err := json.Marshal(rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	m[i] = json_rec
	i++

}

func getByID(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(m[index]))

}

func viewTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	for _, task := range m {
		w.Write(task)
	}
}

func completeTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	var json_slice Record
	err = json.Unmarshal(m[index], &json_slice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json_slice.Completed = true

	updated_json, err := json.Marshal(json_slice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	m[index] = updated_json

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	delete(m, index)
}

func main() {
	http.HandleFunc("/", hellohandler)

	http.HandleFunc("POST /task", addTask)
	http.HandleFunc("GET /task/{id}", getByID)
	http.HandleFunc("GET /task", viewTask)
	http.HandleFunc("PUT /task/{id}", completeTask)
	http.HandleFunc("DELETE /task/{id}", deleteTask)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Not able to start server")
	}
}
