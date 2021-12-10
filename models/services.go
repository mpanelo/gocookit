package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Services struct {
	User   UserService
	Recipe RecipeService
	Image  ImageService
	db     *gorm.DB
}

type ServicesConfig func(*Services) error

func WithGorm(connInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(postgres.Open(connInfo), &gorm.Config{})

		if err != nil {
			return err
		}

		s.db = db
		return nil
	}
}

func WithUser(hmacKey, pepper string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, hmacKey, pepper)
		return nil
	}
}

func WithRecipe() ServicesConfig {
	return func(s *Services) error {
		s.Recipe = NewRecipesService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func WithLogMode(enabled bool) ServicesConfig {
	return func(s *Services) error {
		if enabled {
			s.db.Logger = logger.Default.LogMode(logger.Info)
		}
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}

func (s *Services) DestructiveReset() error {
	if err := s.db.Migrator().DropTable(&User{}, &Recipe{}); err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Recipe{})
}

func (s *Services) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
