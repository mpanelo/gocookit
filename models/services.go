package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Services struct {
	User   UserService
	Recipe RecipeService
	Images ImageService
	db     *gorm.DB
}

func NewServices(dsn string) *Services {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return &Services{
		User:   NewUserService(db),
		Recipe: NewRecipesService(db),
		Images: NewImageService(),
		db:     db,
	}
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
