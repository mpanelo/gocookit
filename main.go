package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpanelo/gocookit/controllers"
)

func main() {
	staticCT := controllers.NewStatic()

	// create router
	router := mux.NewRouter()

	// register routes
	router.Handle("/", staticCT.Home)

	// start http server
	http.ListenAndServe(":8000", router)
}
