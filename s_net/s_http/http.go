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
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		dir := http.Dir(home) // makes the user home directory into a file system
		bashrcFile, err := dir.Open(".bashrc")
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		defer bashrcFile.Close()

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("flusher is not implemented")
			w.WriteHeader(500)
			return
		}

		fileScanner := bufio.NewScanner(bashrcFile)
		// fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			_, w_err := w.Write([]byte(fileScanner.Text() + "\n"))
			if w_err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}
			time.Sleep(500 * time.Millisecond)
			flusher.Flush()
		}

		_, w_err := w.Write([]byte(""))
		if w_err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	})

	http.ListenAndServe(":5000", nil)
}
