package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// create router
	router := mux.NewRouter()

	// register routes
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Hi!")
	})

	// start http server
	http.ListenAndServe(":8000", router)
}
