package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/myvideo", func(w http.ResponseWriter, r *http.Request) {
		video, err := os.ReadFile("myvideo.mp4")
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Error reading file:", err)
			return
		}

		w.Write(video)

	})

	http.ListenAndServe(":5000", nil)
}
