package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/myfiles/notes.md", func(w http.ResponseWriter, r *http.Request) {
		// response ETag - for conditional request
		// w.Header().Set("ETag", "girlsLikeYou")

		file, err := os.Open("notes.md")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(500)
			fmt.Fprintln(w, "Internal Server Error")
			return
		}

		fileStat, _ := file.Stat()

		http.ServeContent(w, r, fileStat.Name(), fileStat.ModTime(), file)
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
			fmt.Fprintln(w, "Internal Server Error")
			return
		}

		fmt.Fprintln(w, "Upload success!")
	})

	http.ListenAndServe(":5000", nil)
}