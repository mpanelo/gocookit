package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mpanelo/gocookit/context"
	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/views"
)

const (
	RouteRecipeEdit = "routeRecipeEdit"
)

type Recipes struct {
	NewView  *views.View
	EditView *views.View
	rs       models.RecipeService
	router   *mux.Router
}

func NewRecipes(rs models.RecipeService, router *mux.Router) *Recipes {
	return &Recipes{
		NewView:  views.NewView("recipes/new"),
		EditView: views.NewView("recipes/edit"),
		rs:       rs,
		router:   router,
	}
}

type RecipeForm struct {
	Title string `schema:"title"`
}

func (rc *Recipes) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form RecipeForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlertDanger(err)
		rc.NewView.Render(rw, r, vd)
		return
	}

	user := context.User(r.Context())

	recipe := models.Recipe{
		UserID: user.ID,
		Title:  form.Title,
	}

	err := rc.rs.Create(&recipe)
	if err != nil {
		vd.SetAlertDanger(err)
		rc.NewView.Render(rw, r, vd)
		return
	}

	url, err := rc.router.Get(RouteRecipeEdit).URL("id", fmt.Sprintf("%v", recipe.ID))
	if err != nil {
		http.Redirect(rw, r, "/recipes", http.StatusFound)
		return
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
}

func (rc *Recipes) Edit(rw http.ResponseWriter, r *http.Request) {
	recipe, err := rc.getRecipe(rw, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())

	if recipe.ID != user.ID {
		http.Error(rw, "Recipe not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = recipe
	rc.EditView.Render(rw, r, vd)
}

func (rc *Recipes) getRecipe(rw http.ResponseWriter, r *http.Request) (*models.Recipe, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	recipeID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid recipe ID", http.StatusNotFound)
		return nil, err
	}

	recipe, err := rc.rs.ByID(uint(recipeID))
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(rw, "Recipe not found", http.StatusNotFound)
			return nil, err
		}

		http.Error(rw, "Unable to find recipe", http.StatusInternalServerError)
		return nil, err
	}

	return recipe, nil
}
