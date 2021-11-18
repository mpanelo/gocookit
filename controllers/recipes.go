package controllers

import (
	"github.com/mpanelo/gocookit/models"
	"github.com/mpanelo/gocookit/views"
)

type Recipes struct {
	NewView *views.View
	gs      models.RecipeService
}

func NewRecipes(gs models.RecipeService) *Recipes {
	return &Recipes{
		NewView: views.NewView("recipes/new"),
		gs:      gs,
	}
}
