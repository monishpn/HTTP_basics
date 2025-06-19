package main

import "net/http"

func hellohandler1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/", hellohandler1)
	http.ListenAndServe(":8080", nil)
}
