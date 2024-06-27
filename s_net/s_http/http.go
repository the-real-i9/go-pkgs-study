package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		var lim int64 = 500 * (1024 * 1000)

		r.Body = http.MaxBytesReader(w, r.Body, lim)

		defer r.Body.Close()

		pmferr := r.ParseMultipartForm(lim)
		if pmferr != nil {
			http.Error(w, "File too large.", http.StatusRequestEntityTooLarge)
			fmt.Println("ParseMultipartForm:", pmferr)
			return
		}

		filehs := r.MultipartForm.File["pic"]
		for _, fileh := range filehs {
			fmt.Println(fileh.Filename)

		}

		for _, filehs := range r.MultipartForm.File {
			for _, fileh := range filehs {
				fmt.Println(fileh.Filename)
			}
		}

	})

	http.ListenAndServe(":5000", nil)
}
