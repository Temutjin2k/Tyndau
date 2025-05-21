package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve all files from web/ directory
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	http.ListenAndServe(":7070", mux)
}
