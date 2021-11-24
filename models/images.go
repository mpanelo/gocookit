package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(uint, io.Reader, string) error
	ByRecipeID(uint) ([]string, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) ByRecipeID(recipeID uint) ([]string, error) {
	images, err := filepath.Glob(is.imageDir(recipeID) + "/*")
	if err != nil {
		return nil, err
	}

	for i := range images {
		images[i] = "/" + images[i]
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
