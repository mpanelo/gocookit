package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	RecipeID uint
	Filename string
}

func (i *Image) Path() string {
	return fmt.Sprintf("/images/recipes/%v/%v", i.RecipeID, i.Filename)
}

func (i *Image) RelativePath() string {
	return i.Path()[1:]
}

type ImageService interface {
	Create(uint, io.Reader, string) error
	ByRecipeID(uint) ([]Image, error)
	Delete(*Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) ByRecipeID(recipeID uint) ([]Image, error) {
	pathPrefix := is.imageDir(recipeID)
	imagePaths, err := filepath.Glob(pathPrefix + "/*")
	if err != nil {
		return nil, err
	}

	images := make([]Image, len(imagePaths))
	for i := range imagePaths {
		images[i] = Image{
			RecipeID: recipeID,
			Filename: strings.Replace(imagePaths[i], pathPrefix+"/", "", 1),
		}
	}

	return images, nil
}

func (is *imageService) Create(recipeID uint, src io.Reader, fileName string) error {
	dir, err := is.mkImageDir(recipeID)
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func (is *imageService) mkImageDir(recipeID uint) (string, error) {
	dir := is.imageDir(recipeID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func (is *imageService) imageDir(recipeID uint) string {
	return filepath.Join("images", "recipes", fmt.Sprintf("%d", recipeID))
}
