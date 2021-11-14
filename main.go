package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpanelo/gocookit/controllers"
	"github.com/mpanelo/gocookit/models"
)

const (
	host   = "localhost"
	user   = "postgres"
	dbname = "gocookit_dev"
	port   = "5432"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, user, dbname, port)
	services := models.NewServices(dsn)
	defer services.Close()

	services.AutoMigrate()

	staticCT := controllers.NewStatic()
	usersCT := controllers.NewUsers(services.User)

	router := mux.NewRouter()

	router.Handle("/", staticCT.Home)

	router.Handle("/signup", usersCT.SignUp).Methods(http.MethodGet)
	router.HandleFunc("/users", usersCT.Create).Methods(http.MethodPost)

	http.ListenAndServe(":8000", router)
}
