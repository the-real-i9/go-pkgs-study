package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		// res := "Bro! How ya doin?\n"
		r.Write(w)
		// w.Write([]byte(res))
	})

	http.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		limBody := http.MaxBytesReader(w, r.Body, 10)

		defer limBody.Close()

		data, err := io.ReadAll(limBody)

		fmt.Printf("%s\n", data)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error() + "\n"))
			return
		}

		w.Write([]byte("Post success!\n"))
	})

	http.ListenAndServe(":5000", nil)
}
