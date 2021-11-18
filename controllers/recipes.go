package controllers

import (
	"fmt"
	"net/http"

	"github.com/mpanelo/gocookit/context"
	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/views"
)

type Recipes struct {
	NewView *views.View
	rs      models.RecipeService
}

func NewRecipes(rs models.RecipeService) *Recipes {
	return &Recipes{
		NewView: views.NewView("recipes/new"),
		rs:      rs,
	}
}

type RecipeForm struct {
	Title        string `schema:"title"`
	Description  string `schema:"description"`
	Ingredients  string `schema:"ingredients"`
	Instructions string `schema:"instructions"`
}

func (rc *Recipes) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form RecipeForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlertDanger(err)
		rc.NewView.Render(rw, vd)
		return
	}

	user := context.User(r.Context())

	recipe := models.Recipe{
		UserID:       user.ID,
		Title:        form.Title,
		Description:  form.Description,
		Ingredients:  form.Ingredients,
		Instructions: form.Instructions,
	}

	err := rc.rs.Create(&recipe)
	if err != nil {
		vd.SetAlertDanger(err)
		rc.NewView.Render(rw, vd)
		return
	}

	// TODO redirect to /recipes page
	fmt.Fprintln(rw, recipe)
}
