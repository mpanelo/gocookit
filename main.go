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
	recipesCT := controllers.NewRecipes(services.Recipe)

	router := mux.NewRouter()

	router.Handle("/", staticCT.Home)

	setUsersRoutes(router, usersCT)
	setRecipesRoutes(router, recipesCT)

	http.ListenAndServe(":8000", router)
}

func setUsersRoutes(router *mux.Router, usersCT *controllers.Users) {
	router.Handle("/signup", usersCT.SignUpView).Methods(http.MethodGet)
	router.Handle("/signin", usersCT.SignInView).Methods(http.MethodGet)
	router.HandleFunc("/users", usersCT.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", usersCT.SignIn).Methods(http.MethodPost)

	router.HandleFunc("/whoami", usersCT.Whoami).Methods(http.MethodGet) // TODO delete temporary endpoint
}

func setRecipesRoutes(router *mux.Router, recipesCT *controllers.Recipes) {
	router.Handle("/recipes/new", recipesCT.NewView).Methods(http.MethodGet)
}
