package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpanelo/gocookit/controllers"
	"github.com/mpanelo/gocookit/middleware"
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

	router := mux.NewRouter()

	staticCT := controllers.NewStatic()
	usersCT := controllers.NewUsers(services.User)
	recipesCT := controllers.NewRecipes(services.Recipe, services.Images, router)

	router.Handle("/", staticCT.Home)

	imagesHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))
	router.PathPrefix("/images/").Handler(imagesHandler)

	setUsersRoutes(router, usersCT)
	setRecipesRoutes(router, recipesCT)

	userMw := middleware.User{UserService: services.User}
	http.ListenAndServe(":8000", userMw.Apply(router))
}

func setUsersRoutes(router *mux.Router, usersCT *controllers.Users) {
	router.Handle("/signup", usersCT.SignUpView).Methods(http.MethodGet)
	router.Handle("/signin", usersCT.SignInView).Methods(http.MethodGet)
	router.HandleFunc("/users", usersCT.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", usersCT.SignIn).Methods(http.MethodPost)
}

func setRecipesRoutes(router *mux.Router, recipesCT *controllers.Recipes) {
	requireUserMw := middleware.RequireUser{}
	router.
		Handle("/recipes", requireUserMw.ApplyFn(recipesCT.Index)).
		Methods(http.MethodGet)
	router.
		Handle("/recipes", requireUserMw.ApplyFn(recipesCT.Create)).
		Methods(http.MethodPost)
	router.
		Handle("/recipes/new", requireUserMw.Apply(recipesCT.NewView)).
		Methods(http.MethodGet)
	router.
		Handle("/recipes/{id:[0-9]+}", requireUserMw.ApplyFn(recipesCT.Show)).
		Methods(http.MethodGet).
		Name(controllers.RouteRecipeShow)
	router.
		Handle("/recipes/{id:[0-9]+}/edit", requireUserMw.ApplyFn(recipesCT.Edit)).
		Methods(http.MethodGet).
		Name(controllers.RouteRecipeEdit)
	router.
		Handle("/recipes/{id:[0-9]+}", requireUserMw.ApplyFn(recipesCT.Update)).
		Methods(http.MethodPost)
	router.
		Handle("/recipes/{id:[0-9]+}/images", requireUserMw.ApplyFn(recipesCT.ImageUpload)).
		Methods(http.MethodPost)
	router.
		Handle("/recipes/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(recipesCT.ImageDelete)).
		Methods(http.MethodPost)
}
