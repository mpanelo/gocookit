package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	Title        string `gorm:"not null"`
	Description  string
	Ingredients  string
	Instructions string
	Images       []Image `gorm:"-"`
}

func (r *Recipe) ImagesSplitN(n int) [][]Image {
	buckets := make([][]Image, n)

	for i, image := range r.Images {
		bucketIdx := i % n
		buckets[bucketIdx] = append(buckets[bucketIdx], image)
	}

	return buckets
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
	ByID(uint) (*Recipe, error)
	ByUserID(uint) ([]Recipe, error)
	Create(*Recipe) error
	Update(*Recipe) error
}

type recipeValidator struct {
	RecipeDB
}

func (rv *recipeValidator) Create(recipe *Recipe) error {
	err := runRecipeValidatorFuncs(recipe,
		userIDRequired,
		titleRequired)
	if err != nil {
		return err
	}

	return rv.RecipeDB.Create(recipe)
}

func (rv *recipeValidator) Update(recipe *Recipe) error {
	err := runRecipeValidatorFuncs(recipe,
		userIDRequired,
		titleRequired)
	if err != nil {
		return err
	}

	return rv.RecipeDB.Update(recipe)
}

func userIDRequired(recipe *Recipe) error {
	if recipe.UserID <= 0 {
		return ErrRecipeUserIDRequired
	}
	return nil
}

func titleRequired(recipe *Recipe) error {
	if recipe.Title == "" {
		return ErrRecipeTitleRequired
	}
	return nil
}

type recipeGorm struct {
	db *gorm.DB
}

func (rg *recipeGorm) ByID(id uint) (*Recipe, error) {
	var recipe Recipe
	tx := rg.db.Where("id = ?", id)

	if err := first(tx, &recipe); err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (rg *recipeGorm) ByUserID(userID uint) ([]Recipe, error) {
	var recipes []Recipe
	result := rg.db.Where("user_id", userID).Find(&recipes)
	if result.Error != nil {
		return nil, result.Error
	}
	return recipes, nil
}

func (rg *recipeGorm) Create(recipe *Recipe) error {
	result := rg.db.Create(recipe)
	return result.Error
}

func (rg *recipeGorm) Update(recipe *Recipe) error {
	result := rg.db.Save(recipe)
	return result.Error
}

type recipeValidatorFunc func(*Recipe) error

func runRecipeValidatorFuncs(recipe *Recipe, funcs ...recipeValidatorFunc) error {
	for _, f := range funcs {
		if err := f(recipe); err != nil {
			return err
		}
	}
	return nil
}
