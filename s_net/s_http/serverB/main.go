package serverB

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/redir", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "You we're redirected to /redir")
	})

	http.HandleFunc("/dosmth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "You we're redirected to /dosmth")
	})
	http.ListenAndServe(":5001", nil)
}
