package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		notesFile, err := os.Open("notes.md.html")
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		defer notesFile.Close()

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("flusher is not implemented")
			w.WriteHeader(500)
			return
		}

		fileScanner := bufio.NewScanner(notesFile)

		fileScanner.Split(bufio.ScanBytes)

		for fileScanner.Scan() {
			w.Header().Set("Content-Type", "text/html")
			_, w_err := w.Write(fileScanner.Bytes())
			if w_err != nil {
				log.Println(err)
				return
			}
			// w.Write([]byte("\n"))
			time.Sleep(10 * time.Millisecond)
			flusher.Flush()
		}
	})

	http.ListenAndServe(":5000", nil)
}
