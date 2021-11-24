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
	RouteRecipeShow = "routeRecipeShow"

	maxMultipartFormMemory = 5 << 20 // 5 megabytes
)

type Recipes struct {
	NewView   *views.View
	EditView  *views.View
	IndexView *views.View
	ShowView  *views.View
	rs        models.RecipeService
	is        models.ImageService
	router    *mux.Router
}

func NewRecipes(rs models.RecipeService, is models.ImageService, router *mux.Router) *Recipes {
	return &Recipes{
		NewView:   views.NewView("recipes/new"),
		EditView:  views.NewView("recipes/edit"),
		IndexView: views.NewView("recipes/index"),
		ShowView:  views.NewView("recipes/show"),
		rs:        rs,
		is:        is,
		router:    router,
	}
}

func (rc *Recipes) UploadImages(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data

	recipe, err := rc.getRecipe(rw, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if recipe.UserID != user.ID {
		http.Error(rw, "Recipe not found", http.StatusNotFound)
		return
	}

	vd.Yield = recipe

	err = r.ParseMultipartForm(maxMultipartFormMemory)
	if err != nil {
		vd.SetAlertDanger(err)
		rc.EditView.Render(rw, r, vd)
		return
	}

	for _, imageFile := range r.MultipartForm.File["images"] {
		srcFile, err := imageFile.Open()
		if err != nil {
			vd.SetAlertDanger(err)
			rc.EditView.Render(rw, r, vd)
			return
		}
		defer srcFile.Close()

		err = rc.is.Create(recipe.ID, srcFile, imageFile.Filename)
		if err != nil {
			vd.SetAlertDanger(err)
			rc.EditView.Render(rw, r, vd)
			return
		}
	}

	vd.SetSuccess("Images uploaded successfully")
	rc.EditView.Render(rw, r, vd)
}

func (rc *Recipes) Show(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data

	recipe, err := rc.getRecipe(rw, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if recipe.UserID != user.ID {
		http.Error(rw, "Recipe not found", http.StatusNotFound)
		return
	}

	vd.Yield = recipe
	rc.ShowView.Render(rw, r, vd)
}

func (rc *Recipes) Index(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data

	user := context.User(r.Context())

	recipes, err := rc.rs.ByUserID(user.ID)
	if err != nil {
		vd.SetAlertDanger(err)
		rc.IndexView.Render(rw, r, vd)
		return
	}

	vd.Yield = recipes
	rc.IndexView.Render(rw, r, vd)
}

type RecipeCreateForm struct {
	Title string `schema:"title"`
}

func (rc *Recipes) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form RecipeCreateForm

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
	if recipe.UserID != user.ID {
		http.Error(rw, "Recipe not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = recipe
	rc.EditView.Render(rw, r, vd)
}

type RecipeUpdateForm struct {
	Title        string `schema:"title"`
	Description  string `schema:"description"`
	Ingredients  string `schema:"ingredients"`
	Instructions string `schema:"instructions"`
}

func (rc *Recipes) Update(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form RecipeUpdateForm

	recipe, err := rc.getRecipe(rw, r)
	if err != nil {
		return
	}

	if err := parseForm(r, &form); err != nil {
		vd.SetAlertDanger(err)
		rc.EditView.Render(rw, r, vd)
		return
	}

	user := context.User(r.Context())
	if recipe.UserID != user.ID {
		http.Error(rw, "Recipe not found", http.StatusNotFound)
		return
	}

	recipe.Title = form.Title
	recipe.Description = form.Description
	recipe.Ingredients = form.Ingredients
	recipe.Instructions = form.Instructions

	err = rc.rs.Update(recipe)
	if err != nil {
		vd.SetAlertDanger(err)
		rc.EditView.Render(rw, r, vd)
		return
	}

	url, err := rc.router.Get(RouteRecipeShow).URL("id", fmt.Sprintf("%v", recipe.ID))
	if err != nil {
		http.Redirect(rw, r, "/recipes", http.StatusFound)
		return
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
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

	images, err := rc.is.ByRecipeID(recipe.ID)
	if err != nil {
		http.Error(rw, "Failed to fetch recipe images", http.StatusInternalServerError)
		return nil, err
	}
	recipe.Images = images
	return recipe, nil
}
