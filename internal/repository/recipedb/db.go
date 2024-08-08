package recipedb

import (
	"crud/internal/domain"
)

type DB interface {
	Get(id string) (*domain.Recipe, error)
	Set(id string, recipe *domain.Recipe) error
	Delete(id string) error
	GetAll() (map[string]*domain.Recipe, error) // Метод для получения всех рецептов
}
