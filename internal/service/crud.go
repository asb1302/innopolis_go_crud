package service

import (
	"crud/internal/domain"
	"crud/internal/repository/recipedb"
	"errors"
	"github.com/google/uuid"
	"sort"
)

var recipes recipedb.DB

func Init(DB recipedb.DB) {
	recipes = DB
}

func Get(id string) (*domain.Recipe, error) {
	return recipes.Get(id)
}

func Delete(id string) error {
	return recipes.Delete(id)
}

func AddOrUpd(r *domain.Recipe) error {

	if r.ID == "" {
		r.ID = uuid.New().String()
	}

	return recipes.Set(r.ID, r)
}

func GetPaginated(page, limit int) ([]domain.Recipe, error) {
	allRecipes, err := recipes.GetAll()
	if err != nil {
		return nil, err
	}

	if page < 1 || limit < 1 {
		return nil, errors.New("page и limit должны быть больше 0")
	}

	// Сортируем рецепты по ID, чтобы упорядочить их для последующей пагинации
	recipeIDs := make([]string, 0, len(allRecipes))
	for id := range allRecipes {
		recipeIDs = append(recipeIDs, id)
	}
	sort.Strings(recipeIDs)

	// Определяем начальный и конечный индексы для слайса
	startIndex := (page - 1) * limit
	if startIndex >= len(recipeIDs) {
		return nil, nil // Если начальный индекс за пределами доступных данных, возвращаем пустой срез
	}
	endIndex := startIndex + limit
	if endIndex > len(recipeIDs) {
		endIndex = len(recipeIDs)
	}

	// Получаем слайс ID для текущей страницы
	pageIDs := recipeIDs[startIndex:endIndex]

	// Формируем слайс рецептов для текущей страницы
	pageRecipes := make([]domain.Recipe, len(pageIDs))
	for i, id := range pageIDs {
		pageRecipes[i] = *allRecipes[id]
	}

	return pageRecipes, nil
}

func GetRecipeCount() (int, error) {
	allRecipes, err := recipes.GetAll()
	if err != nil {
		return 0, err
	}
	return len(allRecipes), nil
}
