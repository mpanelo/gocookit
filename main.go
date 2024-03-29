package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/mpanelo/gocookit/controllers"
	"github.com/mpanelo/gocookit/middleware"
	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/rand"
)

func main() {
	isProdFlag := flag.Bool("prod", false, "Provide this flag in production to ensure that a .config file is provided before the application is started")
	flag.Parse()

	cfg := LoadConfig(*isProdFlag)
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.HMACKey, cfg.Pepper),
		models.WithRecipe(),
		models.WithImage(),
	)
	must(err)

	defer services.Close()
	services.AutoMigrate()

	router := mux.NewRouter()

	staticCT := controllers.NewStatic()
	usersCT := controllers.NewUsers(services.User)
	recipesCT := controllers.NewRecipes(services.Recipe, services.Image, router)

	router.Handle("/", staticCT.Home)

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))
	router.PathPrefix("/assets/").Handler(assetsHandler)

	imagesHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("./images/")))
	router.PathPrefix("/images/").Handler(imagesHandler)

	setUsersRoutes(router, usersCT)
	setRecipesRoutes(router, recipesCT)

	b, err := rand.Bytes(32)
	must(err)

	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	userMw := middleware.User{UserService: services.User}

	portStr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Starting gocookit on %s...\n", portStr)
	http.ListenAndServe(portStr, csrfMw(userMw.Apply(router)))
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
