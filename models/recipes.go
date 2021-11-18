package models

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Title  string `gorm:"not null"`
}

type RecipeService interface {
	RecipeDB
}

type recipeService struct {
	RecipeDB
}

func NewRecipesService(db *gorm.DB) RecipeService {
	return &recipeService{&recipeValidator{&recipeGorm{db}}}
}

type RecipeDB interface {
	Create(*Recipe) error
}

type recipeValidator struct {
	RecipeDB
}

type recipeGorm struct {
	db *gorm.DB
}

func (gg *recipeGorm) Create(recipe *Recipe) error {
	result := gg.db.Create(recipe)
	return result.Error
}
