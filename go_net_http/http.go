package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.RedirectHandler("http://localhost:5001"+r.RequestURI, http.StatusFound).ServeHTTP(w, r)
	})

	/* http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	}) */

	http.ListenAndServe(":5000", nil)
}
