package main

import (
	"errors"
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
		var fileSizeLimit int64 = 10

		limBody := http.MaxBytesReader(w, r.Body, fileSizeLimit)

		defer limBody.Close()

		data, err := io.ReadAll(limBody)

		// incomplete read data
		fmt.Printf("%s\n", data)

		if err != nil {
			var mbe *http.MaxBytesError

			if errors.As(err, &mbe) {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				fmt.Fprintf(w, "File too large. Limit is %d bytes.\n", fileSizeLimit)
				return
			}
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintln(w, "Upload success!")
	})

	http.ListenAndServe(":5000", nil)
}
