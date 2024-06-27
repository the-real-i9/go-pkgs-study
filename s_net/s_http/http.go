package main

import (
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
)

func main() {

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		lim := int64(1024 * (math.Pow10(1)))

		r.Body = http.MaxBytesReader(w, r.Body, lim)

		defer r.Body.Close()

		pmferr := r.ParseMultipartForm(1024 * lim)
		if pmferr != nil {
			http.Error(w, "File too large.", http.StatusRequestEntityTooLarge)
			fmt.Println("ParseMultipartForm:", pmferr)
			return
		}

		files := r.MultipartForm.File["pic"]
		file, _ := files[0].Open()
		data, _ := io.ReadAll(file)

		for key, files := range r.MultipartForm.File {
			file, _ := files[0].Open()
			data, _ := io.ReadAll(file)
		}

	})

	http.ListenAndServe(":5000", nil)
}
