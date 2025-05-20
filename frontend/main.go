package main

import (
	"net/http"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", ServeHome)

	http.ListenAndServe(":7070", mux)
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "index.html"))
}
